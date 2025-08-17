package stockmovementvehicleretail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/jihanlugas/warehouse/app/retail"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovementvehicle"
	"github.com/jihanlugas/warehouse/app/vehicle"
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
	Page(loginUser jwt.UserLogin, req request.PageStockmovementvehicleRetail) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateStockmovementvehicleRetail) (err error)
	Update(loginUser jwt.UserLogin, id string, req request.UpdateStockmovementvehicleRetail) (err error)
	Delete(loginUser jwt.UserLogin, id string) (err error)
	SetComplete(loginUser jwt.UserLogin, id string) (err error)
	GenerateDeliveryOrder(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error)
}

type usecase struct {
	retailRepository               retail.Repository
	stockmovementvehicleRepository stockmovementvehicle.Repository
	warehouseRepository            warehouse.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
	vehicleRepository              vehicle.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageStockmovementvehicleRetail) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	pageReq := request.PageStockmovementvehicle{
		Paging:                     req.Paging,
		FromWarehouseID:            loginUser.WarehouseID,
		ProductID:                  req.ProductID,
		VehicleID:                  req.VehicleID,
		RelatedID:                  req.RetailID,
		StockmovementvehicleType:   model.StockmovementvehicleTypeRetail,
		Notes:                      req.Notes,
		StockmovementvehicleStatus: req.StockmovementvehicleStatus,
		StartSentGrossQuantity:     req.StartSentGrossQuantity,
		StartSentTareQuantity:      req.StartSentTareQuantity,
		StartSentNetQuantity:       req.StartSentNetQuantity,
		StartSentTime:              req.StartSentTime,
		EndSentGrossQuantity:       req.EndSentGrossQuantity,
		EndSentTareQuantity:        req.EndSentTareQuantity,
		EndSentNetQuantity:         req.EndSentNetQuantity,
		EndSentTime:                req.EndSentTime,
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

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
		return vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	}

	return vStockmovementvehicle, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateStockmovementvehicleRetail) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle
	var vWarehouse model.WarehouseView
	var tVehicle model.Vehicle
	var vRetail model.RetailView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, loginUser.WarehouseID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.warehouseRepository.Name(), err))
	}

	if !vWarehouse.IsRetail {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create retail"))
	}

	vRetail, err = u.retailRepository.GetViewById(conn, req.RetailID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.retailRepository.Name(), err))
	}

	if vRetail.RetailStatus != model.RetailStatusOpen {
		return errors.New(fmt.Sprintf("unable to create retail with status %s", strings.ToLower(string(vRetail.RetailStatus))))
	}

	tx := conn.Begin()

	if req.IsNewVehiclerdriver {
		tVehicle = model.Vehicle{
			ID:          utils.GetUniqueID(),
			WarehouseID: loginUser.WarehouseID,
			PlateNumber: strings.ToUpper(req.PlateNumber),
			Name:        req.VehicleName,
			NIK:         req.NIK,
			DriverName:  req.DriverName,
			PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.vehicleRepository.Create(tx, tVehicle)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.vehicleRepository.Name(), err))
		}
	} else {
		tVehicle, err = u.vehicleRepository.GetTableById(conn, req.VehicleID)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.vehicleRepository.Name(), err))
		}
	}

	tStockmovementvehicle = model.Stockmovementvehicle{
		FromLocationID:             vWarehouse.LocationID,
		FromWarehouseID:            vWarehouse.ID,
		RelatedID:                  req.RetailID,
		ProductID:                  req.ProductID,
		VehicleID:                  tVehicle.ID,
		StockmovementvehicleType:   model.StockmovementvehicleTypeRetail,
		Notes:                      req.Notes,
		SentGrossQuantity:          req.SentGrossQuantity,
		SentTareQuantity:           req.SentTareQuantity,
		SentNetQuantity:            req.SentNetQuantity,
		StockmovementvehicleStatus: model.StockmovementvehicleStatusLoading,
		CreateBy:                   loginUser.UserID,
		UpdateBy:                   loginUser.UserID,
	}
	err = u.stockmovementvehicleRepository.Create(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateStockmovementvehicleRetail) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusLoading {
		return errors.New(fmt.Sprintf("unable to update data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	tx := conn.Begin()

	tStockmovementvehicle.SentTareQuantity = req.SentTareQuantity
	tStockmovementvehicle.SentGrossQuantity = req.SentGrossQuantity
	tStockmovementvehicle.SentNetQuantity = req.SentNetQuantity
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

func (u usecase) Delete(loginUser jwt.UserLogin, id string) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusLoading {
		return errors.New(fmt.Sprintf("unable to delete data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	tx := conn.Begin()

	err = u.stockmovementvehicleRepository.Delete(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.stockmovementvehicleRepository.Name(), err))
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

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.FromWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusLoading {
		return errors.New(fmt.Sprintf("unable to update data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	now := time.Now()
	tx := conn.Begin()

	tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, tStockmovementvehicle.FromWarehouseID, tStockmovementvehicle.ProductID)
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

	tStock.Quantity = tStock.Quantity - tStockmovementvehicle.SentNetQuantity
	err = u.stockRepository.Save(tx, tStock)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockRepository.Name(), err))
	}

	tStockmovementvehicle.SentTime = &now
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
		StocklogType:           model.StocklogTypeOut,
		GrossQuantity:          tStockmovementvehicle.SentGrossQuantity,
		TareQuantity:           tStockmovementvehicle.SentTareQuantity,
		NetQuantity:            tStockmovementvehicle.SentNetQuantity,
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

func (u usecase) GenerateDeliveryOrder(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vStockmovementvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id, "Vehicle", "FromWarehouse", "Product", "Retail", "Retail.Customer")
	if err != nil {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if vStockmovementvehicle.FromWarehouse == nil {
		log.Info("warehouse not found")
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorDataNotFound)
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
		return pdfBytes, vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	}

	if vStockmovementvehicle.StockmovementvehicleStatus != model.StockmovementvehicleStatusCompleted {
		return pdfBytes, vStockmovementvehicle, errors.New(fmt.Sprintf("unable to generate data with status %s", strings.ToLower(string(vStockmovementvehicle.StockmovementvehicleStatus))))
	}

	pdfBytes, err = u.generateDeliveryOrder(vStockmovementvehicle)

	return pdfBytes, vStockmovementvehicle, err
}

func (u usecase) generateDeliveryOrder(vStockmovementvehicle model.StockmovementvehicleView) (pdfBytes []byte, err error) {
	tmpl := template.New("retail-delivery-order.html").Funcs(template.FuncMap{
		"displayImagePhotoId": utils.GetPhotoUrlById,
		"displayDate":         utils.DisplayDate,
		"displayDatetime":     utils.DisplayDatetime,
		"displayNumber":       utils.DisplayNumber,
		"displayMoney":        utils.DisplayMoney,
		"displayPhoneNumber":  utils.DisplayPhoneNumber,
		"displaySpkNumber":    utils.DisplaySpkNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/retail-delivery-order.html")
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

func NewUsecase(retailRepository retail.Repository, stockmovementvehicleRepository stockmovementvehicle.Repository, warehouseRepository warehouse.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, vehicleRepository vehicle.Repository) Usecase {
	return &usecase{
		retailRepository:               retailRepository,
		stockmovementvehicleRepository: stockmovementvehicleRepository,
		warehouseRepository:            warehouseRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
		vehicleRepository:              vehicleRepository,
	}
}
