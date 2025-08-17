package transferin

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovementvehicle"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageTransferin) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	Update(loginUser jwt.UserLogin, id string, req request.UpdateTransferin) (err error)
	SetUnloading(loginUser jwt.UserLogin, id string) (err error)
	SetComplete(loginUser jwt.UserLogin, id string) (err error)
	GenerateDeliveryRecipt(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error)
}

type usecase struct {
	stockmovementvehicleRepository stockmovementvehicle.Repository
	warehouseRepository            warehouse.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageTransferin) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	pageReq := request.PageStockmovementvehicle{
		Paging:                     req.Paging,
		ToWarehouseID:              loginUser.WarehouseID,
		ProductID:                  req.ProductID,
		VehicleID:                  req.VehicleID,
		StockmovementvehicleType:   model.StockmovementvehicleTypeTransfer,
		Notes:                      req.Notes,
		StockmovementvehicleStatus: req.StockmovementvehicleStatus,
		StartReceivedGrossQuantity: req.StartReceivedGrossQuantity,
		StartReceivedTareQuantity:  req.StartReceivedTareQuantity,
		StartReceivedNetQuantity:   req.StartReceivedNetQuantity,
		StartReceivedTime:          req.StartReceivedTime,
		EndReceivedGrossQuantity:   req.EndReceivedGrossQuantity,
		EndReceivedTareQuantity:    req.EndReceivedTareQuantity,
		EndReceivedNetQuantity:     req.EndReceivedNetQuantity,
		EndReceivedTime:            req.EndReceivedTime,
		CreateName:                 req.CreateName,
		StartCreateDt:              req.StartCreateDt,
		EndCreateDt:                req.EndCreateDt,
		Preloads:                   req.Preloads,
	}

	vStockmovementvehicles, count, err = u.stockmovementvehicleRepository.Page(conn, pageReq)
	if err != nil {
		return vStockmovementvehicles, count, err
	}

	return vStockmovementvehicles, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.ToWarehouseID) {
		return vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	}

	return vStockmovementvehicle, err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateTransferin) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.ToWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusUnloading {
		return errors.New(fmt.Sprintf("unable to update data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	tx := conn.Begin()

	tStockmovementvehicle.ReceivedTareQuantity = req.ReceivedTareQuantity
	tStockmovementvehicle.ReceivedGrossQuantity = req.ReceivedGrossQuantity
	tStockmovementvehicle.ReceivedNetQuantity = req.ReceivedNetQuantity
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetUnloading(loginUser jwt.UserLogin, id string) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.ToWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusInTransit {
		return errors.New(fmt.Sprintf("unable to update data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	now := time.Now()
	tx := conn.Begin()

	tStockmovementvehicle.StockmovementvehicleStatus = model.StockmovementvehicleStatusUnloading
	tStockmovementvehicle.ReceivedTime = &now
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetComplete(loginUser jwt.UserLogin, id string) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle
	var tStock model.Stock
	var tStocklog model.Stocklog

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.ToWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusUnloading {
		return errors.New(fmt.Sprintf("unable to update data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	tx := conn.Begin()

	tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, tStockmovementvehicle.ToWarehouseID, tStockmovementvehicle.ProductID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockRepository.Name(), err))
		}
		tStock = model.Stock{
			ID:          utils.GetUniqueID(),
			WarehouseID: loginUser.WarehouseID,
			ProductID:   tStockmovementvehicle.ProductID,
			Quantity:    0,
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.stockRepository.Create(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
		}
	}

	tStock.Quantity = tStock.Quantity + tStockmovementvehicle.ReceivedNetQuantity
	err = u.stockRepository.Save(tx, tStock)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockRepository.Name(), err))
	}

	tStockmovementvehicle.Shrinkage = tStockmovementvehicle.SentNetQuantity - tStockmovementvehicle.ReceivedNetQuantity
	tStockmovementvehicle.StockmovementvehicleStatus = model.StockmovementvehicleStatusCompleted
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	tStocklog = model.Stocklog{
		WarehouseID:            loginUser.WarehouseID,
		StockID:                tStock.ID,
		StockmovementvehicleID: tStockmovementvehicle.ID,
		ProductID:              tStockmovementvehicle.ProductID,
		VehicleID:              tStockmovementvehicle.VehicleID,
		StocklogType:           model.StocklogTypeIn,
		GrossQuantity:          tStockmovementvehicle.ReceivedGrossQuantity,
		TareQuantity:           tStockmovementvehicle.ReceivedTareQuantity,
		NetQuantity:            tStockmovementvehicle.ReceivedNetQuantity,
		CurrentQuantity:        tStock.Quantity,
		CreateBy:               loginUser.UserID,
		UpdateBy:               loginUser.UserID,
	}
	err = u.stocklogRepository.Create(tx, tStocklog)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.stocklogRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) GenerateDeliveryRecipt(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id, "Vehicle", "FromWarehouse", "ToWarehouse", "Product")
	if err != nil {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if vStockmovementvehicle.ToWarehouse == nil {
		log.Info("warehouse not found")
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorDataNotFound)
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.ToWarehouseID) {
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	}

	if vStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusCompleted {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("unable to generate data with status %s", strings.ToLower(string(vStockmovementvehicle.StockmovementvehicleStatus))))
	}

	pdfBytes, err = u.generateDeliveryRecipt(vStockmovementvehicle)

	return pdfBytes, vStockmovementvehicle, err
}

func (u usecase) generateDeliveryRecipt(vStockmovementvehicle model.StockmovementvehicleView) (pdfBytes []byte, err error) {
	tmpl := template.New("delivery-recipt.html").Funcs(template.FuncMap{
		"displayNumberMinus": func(a, b float64) string {
			return utils.DisplayNumber(a - b)
		},
		"displayImagePhotoId": utils.GetPhotoUrlById,
		"displayDate":         utils.DisplayDate,
		"displayDatetime":     utils.DisplayDatetime,
		"displayNumber":       utils.DisplayNumber,
		"displayMoney":        utils.DisplayMoney,
		"displayPhoneNumber":  utils.DisplayPhoneNumber,
		"displaySpkNumber":    utils.DisplaySpkNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/delivery-recipt.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vStockmovementvehicle); err != nil {
		return pdfBytes, err
	}

	// Simpan HTML render ke file sementara
	tempHTMLFile := "temp.html"
	if err := os.WriteFile(tempHTMLFile, buf.Bytes(), 0644); err != nil {
		return pdfBytes, err
	}
	defer os.Remove(tempHTMLFile)

	return utils.GeneratePDFWithChromedp(tempHTMLFile)
}

func NewUsecase(stockmovementvehicleRepository stockmovementvehicle.Repository, warehouseRepository warehouse.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository) Usecase {
	return &usecase{
		stockmovementvehicleRepository: stockmovementvehicleRepository,
		warehouseRepository:            warehouseRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
	}
}
