package outbound

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
	"strings"
	"time"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageOutbound) (vOutbounds []model.OutboundView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOutbound model.OutboundView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateOutbound) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateOutbound) error
	SetSent(loginUser jwt.UserLogin, id string) error
	GenerateDeliveryOrder(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vOutbound model.OutboundView, err error)
}

type usecase struct {
	outboundRepository             Repository
	warehouseRepository            warehouse.Repository
	vehicleRepository              vehicle.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
	stockmovementRepository        stockmovement.Repository
	stockmovementvehicleRepository stockmovementvehicle.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageOutbound) (vOutbounds []model.OutboundView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOutbounds, count, err = u.outboundRepository.Page(conn, req)
	if err != nil {
		return vOutbounds, count, err
	}

	return vOutbounds, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vOutbound model.OutboundView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOutbound, err = u.outboundRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vOutbound, errors.New(fmt.Sprintf("failed to get %s: %v", u.outboundRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vOutbound.WarehouseID) {
		return vOutbound, errors.New(response.ErrorHandlerIDOR)
	}

	return vOutbound, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateOutbound) error {
	var err error
	var vWarehouse model.WarehouseView
	var tVehicle model.Vehicle
	var tStockmovement model.Stockmovement
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, req.FromWarehouseID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.warehouseRepository.Name(), err))
	}

	if !vWarehouse.IsOutbound {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create outbound"))
	}

	if req.IsNewVehiclerdriver {
		tVehicle = model.Vehicle{
			ID:          utils.GetUniqueID(),
			WarehouseID: req.FromWarehouseID,
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

	tStockmovement = model.Stockmovement{
		ID:              utils.GetUniqueID(),
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		ProductID:       req.ProductID,
		Type:            model.StockmovementTypeTransfer,
		Remark:          req.Remark,
		CreateBy:        loginUser.UserID,
		UpdateBy:        loginUser.UserID,
	}
	err = u.stockmovementRepository.Create(tx, tStockmovement)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockmovementRepository.Name(), err))
	}

	tStockmovementvehicle = model.Stockmovementvehicle{
		StockmovementID:        tStockmovement.ID,
		ProductID:              req.ProductID,
		VehicleID:              tVehicle.ID,
		StockmovementvehicleID: req.StockmovementvehicleID,
		SentGrossQuantity:      req.SentGrossQuantity,
		SentTareQuantity:       req.SentTareQuantity,
		SentNetQuantity:        req.SentNetQuantity,
		Status:                 model.StockmovementvehicleStatusLoading,
		CreateBy:               loginUser.UserID,
		UpdateBy:               loginUser.UserID,
	}
	err = u.stockmovementvehicleRepository.Create(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create stockmovement vehicle: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateOutbound) error {
	var err error
	var vOutbound model.OutboundView
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOutbound, err = u.outboundRepository.GetViewById(conn, id, "Warehouse")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.outboundRepository.Name(), err))
	}

	if vOutbound.Warehouse != nil && !vOutbound.Warehouse.IsOutbound {
		return errors.New(fmt.Sprint("this warehouse is not allowed to update outbound"))
	}

	if vOutbound.SentTime != nil {
		return errors.New("this vehicle already sent")
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vOutbound.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tStockmovementvehicle.SentGrossQuantity = req.SentGrossQuantity
	tStockmovementvehicle.SentTareQuantity = req.SentTareQuantity
	tStockmovementvehicle.SentNetQuantity = req.SentNetQuantity
	tStockmovementvehicle.UpdateBy = loginUser.UserID
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) SetSent(loginUser jwt.UserLogin, id string) error {
	var err error
	var vOutbound model.OutboundView
	var tStock model.Stock
	var tStocklog model.Stocklog
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOutbound, err = u.outboundRepository.GetViewById(conn, id, "Warehouse")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.outboundRepository.Name(), err))
	}

	if vOutbound.Status != model.StockmovementvehicleStatusLoading {
		return errors.New("this vehicle already sent")
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vOutbound.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if tStockmovementvehicle.SentNetQuantity <= 0 && tStockmovementvehicle.SentGrossQuantity <= 0 {
		return errors.New("data gross or net quantity is zero")
	}

	tx := conn.Begin()

	now := time.Now()
	tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, vOutbound.WarehouseID, vOutbound.ProductID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockRepository.Name(), err))
		}
		tStock = model.Stock{
			ID:          utils.GetUniqueID(),
			WarehouseID: vOutbound.WarehouseID,
			ProductID:   vOutbound.ProductID,
			Quantity:    0,
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
		}
	}

	tStockmovementvehicle.SentTime = &now
	tStockmovementvehicle.Status = model.StockmovementvehicleStatusInTransit
	tStockmovementvehicle.UpdateBy = loginUser.UserID
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	// update stock on warehouse if transfer source from warehouse
	if vOutbound.SentNetQuantity != 0 && vOutbound.StockmovementvehicleID == "" {
		CurrentQuantity := tStock.Quantity - vOutbound.SentNetQuantity
		tStock.Quantity = CurrentQuantity
		tStock.UpdateBy = loginUser.UserID
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockRepository.Name(), err))
		}

		tStocklog = model.Stocklog{
			WarehouseID:            vOutbound.WarehouseID,
			StockID:                tStock.ID,
			StockmovementID:        tStockmovementvehicle.StockmovementID,
			StockmovementvehicleID: tStockmovementvehicle.ID,
			ProductID:              tStockmovementvehicle.ProductID,
			VehicleID:              tStockmovementvehicle.VehicleID,
			Type:                   model.StockLogTypeOut,
			GrossQuantity:          vOutbound.SentGrossQuantity,
			TareQuantity:           vOutbound.SentTareQuantity,
			NetQuantity:            vOutbound.SentNetQuantity,
			CurrentQuantity:        CurrentQuantity,
			CreateBy:               loginUser.UserID,
			UpdateBy:               loginUser.UserID,
		}
		err = u.stocklogRepository.Create(tx, tStocklog)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stocklogRepository.Name(), err))
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) GenerateDeliveryOrder(loginUser jwt.UserLogin, id string) (pdfBytes []byte, vOutbound model.OutboundView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vOutbound, err = u.outboundRepository.GetViewById(conn, id, "Warehouse", "Vehicle", "Stockmovement", "Stockmovement.FromWarehouse", "Stockmovement.ToWarehouse", "Product")
	if err != nil {
		return pdfBytes, vOutbound, errors.New(fmt.Sprintf("failed to get %s: %v", u.outboundRepository.Name(), err))
	}

	if vOutbound.Warehouse == nil {
		log.Info("warehouse not found")
		return pdfBytes, vOutbound, errors.New(response.ErrorDataNotFound)
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vOutbound.WarehouseID) {
		return pdfBytes, vOutbound, errors.New(response.ErrorHandlerIDOR)
	}

	pdfBytes, err = u.generateDeliveryOrder(vOutbound)

	return pdfBytes, vOutbound, err
}

func (u usecase) generateDeliveryOrder(vOutbound model.OutboundView) (pdfBytes []byte, err error) {
	tmpl := template.New("delivery-order.html").Funcs(template.FuncMap{
		"displayDate":        utils.DisplayDate,
		"displayDatetime":    utils.DisplayDatetime,
		"displayNumber":      utils.DisplayNumber,
		"displayMoney":       utils.DisplayMoney,
		"displayPhoneNumber": utils.DisplayPhoneNumber,
	})

	// Parse template setelah fungsi didaftarkan
	tmpl, err = tmpl.ParseFiles("assets/template/delivery-order.html")
	if err != nil {
		return pdfBytes, err
	}

	// Render template ke buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vOutbound); err != nil {
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

func NewUsecase(outboundRepository Repository, warehouseRepository warehouse.Repository, vehicleRepository vehicle.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, stockmovementRepository stockmovement.Repository, stockmovementvehicleRepository stockmovementvehicle.Repository) Usecase {
	return &usecase{
		outboundRepository:             outboundRepository,
		warehouseRepository:            warehouseRepository,
		vehicleRepository:              vehicleRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
		stockmovementRepository:        stockmovementRepository,
		stockmovementvehicleRepository: stockmovementvehicleRepository,
	}
}
