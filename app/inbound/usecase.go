package inbound

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovement"
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
	"html/template"
	"os"
	"time"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageInbound) (vInbounds []model.InboundView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vInbound model.InboundView, err error)
	Update(loginUser jwt.UserLogin, id string, req request.UpdateInbound) error
	SetRecived(loginUser jwt.UserLogin, id string) error
	GenerateDeliveryRecipt(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vInbound model.InboundView, err error)
}

type usecase struct {
	inboundRepository              Repository
	warehouseRepository            warehouse.Repository
	vehicleRepository              vehicle.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
	stockmovementRepository        stockmovement.Repository
	stockmovementvehicleRepository stockmovementvehicle.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageInbound) (vInbounds []model.InboundView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vInbounds, count, err = u.inboundRepository.Page(conn, req)
	if err != nil {
		return vInbounds, count, err
	}

	return vInbounds, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vInbound model.InboundView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vInbound, err = u.inboundRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vInbound, errors.New(fmt.Sprint("failed to get inbound: ", err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vInbound.WarehouseID) {
		return vInbound, errors.New(response.ErrorHandlerIDOR)
	}

	return vInbound, err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateInbound) error {
	var err error
	var vInbound model.InboundView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vInbound, err = u.inboundRepository.GetViewById(conn, id, "Warehouse", "Stockmovement")
	if err != nil {
		return errors.New(fmt.Sprint("failed to get inbound: ", err))
	}

	if vInbound.Warehouse != nil && !vInbound.Warehouse.IsInbound {
		return errors.New(fmt.Sprint("this warehouse is not allowed to update inbound"))
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get stockmovementvehicle: ", err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vInbound.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	if tStockmovementvehicle.RecivedTime != nil {
		return errors.New("unable to update data")
	}

	tStockmovementvehicle.RecivedGrossQuantity = req.RecivedGrossQuantity
	tStockmovementvehicle.RecivedTareQuantity = req.RecivedTareQuantity
	tStockmovementvehicle.RecivedNetQuantity = req.RecivedNetQuantity
	tStockmovementvehicle.UpdateBy = loginUser.UserID

	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update stockmovementvehicle: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetRecived(loginUser jwt.UserLogin, id string) error {
	var err error
	var vInbound model.InboundView
	var tStock model.Stock
	var tStocklog model.Stocklog
	var tStockmovementvehicle model.Stockmovementvehicle

	now := time.Now()

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vInbound, err = u.inboundRepository.GetViewById(conn, id, "Warehouse", "Stockmovement")
	if err != nil {
		return errors.New(fmt.Sprint("failed to get inbound: ", err))
	}

	if vInbound.Warehouse != nil && !vInbound.Warehouse.IsInbound {
		return errors.New(fmt.Sprint("this warehouse is not allowed to update inbound"))
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get stockmovementvehicle: ", err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vInbound.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	if tStockmovementvehicle.RecivedTime != nil {
		return errors.New("unable to update data")
	}

	tStockmovementvehicle.UpdateBy = loginUser.UserID

	tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, vInbound.WarehouseID, vInbound.ProductID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("failed to get stock %s: %v", u.stockRepository.Name(), err))
		}
		tStock = model.Stock{
			ID:          utils.GetUniqueID(),
			WarehouseID: vInbound.WarehouseID,
			ProductID:   vInbound.ProductID,
			Quantity:    0,
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
		}
	}

	if tStockmovementvehicle.RecivedNetQuantity == 0 {
		return errors.New(fmt.Sprint("failed to update stockmovementvehicle: weight cannot be zero"))
	}

	if tStockmovementvehicle.RecivedNetQuantity != tStockmovementvehicle.RecivedGrossQuantity-tStockmovementvehicle.RecivedTareQuantity {
		return errors.New(fmt.Sprint("failed to update stockmovementvehicle: weight dosent match"))
	}
	tStockmovementvehicle.RecivedTime = &now

	CurrentQuantity := 0.0
	CurrentQuantity = tStock.Quantity + tStockmovementvehicle.RecivedNetQuantity
	tStock.Quantity = CurrentQuantity
	tStock.UpdateBy = loginUser.UserID
	err = u.stockRepository.Save(tx, tStock)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update stock: ", err))
	}

	tStocklog = model.Stocklog{
		WarehouseID:            vInbound.WarehouseID,
		StockID:                tStock.ID,
		StockmovementID:        tStockmovementvehicle.StockmovementID,
		StockmovementvehicleID: tStockmovementvehicle.ID,
		ProductID:              tStockmovementvehicle.ProductID,
		VehicleID:              tStockmovementvehicle.VehicleID,
		Type:                   model.StockLogTypeIn,
		GrossQuantity:          tStockmovementvehicle.RecivedGrossQuantity,
		TareQuantity:           tStockmovementvehicle.RecivedTareQuantity,
		NetQuantity:            tStockmovementvehicle.RecivedNetQuantity,
		CurrentQuantity:        CurrentQuantity,
		CreateBy:               loginUser.UserID,
		UpdateBy:               loginUser.UserID,
	}
	err = u.stocklogRepository.Create(tx, tStocklog)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.stocklogRepository.Name(), err))
	}

	// process from warehouse data if from warehouse only have gross data
	if tStockmovementvehicle.SentNetQuantity == 0 {
		tStockmovementvehicle.SentTareQuantity = tStockmovementvehicle.RecivedTareQuantity
		tStockmovementvehicle.SentNetQuantity = tStockmovementvehicle.SentGrossQuantity - tStockmovementvehicle.SentTareQuantity

		tFromStock, err := u.stockRepository.GetTableByWarehouseIdAndProductId(tx, vInbound.Stockmovement.FromWarehouseID, vInbound.ProductID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(fmt.Sprint("failed to get from stock: ", err))
			}
			tFromStock = model.Stock{
				ID:          utils.GetUniqueID(),
				WarehouseID: vInbound.Stockmovement.FromWarehouseID,
				ProductID:   vInbound.ProductID,
				Quantity:    0,
				CreateBy:    loginUser.UserID,
				UpdateBy:    loginUser.UserID,
			}
			err = u.stockRepository.Save(tx, tFromStock)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
			}
		}

		FromCurrentQuantity := tFromStock.Quantity - tStockmovementvehicle.SentNetQuantity
		tFromStock.Quantity = FromCurrentQuantity
		tFromStock.UpdateBy = loginUser.UserID
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprint("failed to update from stock: ", err))
		}

		tStocklog = model.Stocklog{
			WarehouseID:            tFromStock.WarehouseID,
			StockID:                tFromStock.ID,
			StockmovementID:        vInbound.StockmovementID,
			StockmovementvehicleID: vInbound.ID,
			ProductID:              vInbound.ProductID,
			VehicleID:              vInbound.VehicleID,
			Type:                   model.StockLogTypeOut,
			GrossQuantity:          tStockmovementvehicle.SentGrossQuantity,
			TareQuantity:           tStockmovementvehicle.SentTareQuantity,
			NetQuantity:            tStockmovementvehicle.SentNetQuantity,
			CurrentQuantity:        FromCurrentQuantity,
			CreateBy:               loginUser.UserID,
			UpdateBy:               loginUser.UserID,
		}
		err = u.stocklogRepository.Create(tx, tStocklog)
		if err != nil {
			return errors.New(fmt.Sprint("failed to create from stocklog: ", err))
		}
	}

	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update stockmovementvehicle: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) GenerateDeliveryRecipt(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vInbound model.InboundView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vInbound, err = u.inboundRepository.GetViewById(conn, id, "Warehouse", "Stockmovement")
	if err != nil {
		return pdfBytes, vInbound, errors.New(fmt.Sprint("failed to get inbound: ", err))
	}

	if vInbound.Warehouse == nil {
		log.Info("warehouse not found")
		return pdfBytes, vInbound, errors.New(response.ErrorDataNotFound)
	}

	if vInbound.RecivedTime == nil {
		return pdfBytes, vInbound, errors.New(fmt.Sprint("failed to get inbound: recived time not found"))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vInbound.WarehouseID) {
		return pdfBytes, vInbound, errors.New(response.ErrorHandlerIDOR)
	}

	pdfBytes, err = u.generateDeliveryRecipt(vInbound)

	return pdfBytes, vInbound, err
}

func (u usecase) generateDeliveryRecipt(vInbound model.InboundView) (pdfBytes []byte, err error) {
	tmpl := template.New("delivery-recipt.html").Funcs(template.FuncMap{
		"displayLembar": func(lembar int64) string {
			return fmt.Sprintf("%s Lembar", utils.DisplayNumber(lembar))
		},
		"displayDuplex": func(isDuplex bool) string {
			if isDuplex {
				return "2 Muka"
			}
			return "1 Muka"
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
	tmpl, err = tmpl.ParseFiles("assets/template/delivery-order.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vInbound); err != nil {
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

func NewUsecase(inboundRepository Repository, warehouseRepository warehouse.Repository, vehicleRepository vehicle.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, stockmovementRepository stockmovement.Repository, stockmovementvehicleRepository stockmovementvehicle.Repository) Usecase {
	return &usecase{
		inboundRepository:              inboundRepository,
		warehouseRepository:            warehouseRepository,
		vehicleRepository:              vehicleRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
		stockmovementRepository:        stockmovementRepository,
		stockmovementvehicleRepository: stockmovementvehicleRepository,
	}
}
