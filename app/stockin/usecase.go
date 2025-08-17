package stockin

import (
	"errors"
	"fmt"
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
	"gorm.io/gorm"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageStockin) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateStockin) (err error)
	Delete(loginUser jwt.UserLogin, id string) (err error)
	SetComplete(loginUser jwt.UserLogin, id string) (err error)
}

type usecase struct {
	stockmovementvehicleRepository stockmovementvehicle.Repository
	warehouseRepository            warehouse.Repository
	stockRepository                stock.Repository
	stocklogRepository             stocklog.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageStockin) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	pageReq := request.PageStockmovementvehicle{
		Paging:                     req.Paging,
		ToWarehouseID:              loginUser.WarehouseID,
		ProductID:                  req.ProductID,
		StockmovementvehicleType:   model.StockmovementvehicleTypeIn,
		Notes:                      req.Notes,
		StockmovementvehicleStatus: req.StockmovementvehicleStatus,
		StartReceivedNetQuantity:   req.StartNetQuantity,
		EndReceivedNetQuantity:     req.EndNetQuantity,
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

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateStockin) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle
	var vWarehouse model.WarehouseView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, loginUser.WarehouseID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.warehouseRepository.Name(), err))
	}

	if !vWarehouse.IsStockin {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create stockin"))
	}

	tx := conn.Begin()
	now := time.Now()

	tStockmovementvehicle = model.Stockmovementvehicle{
		ToLocationID:               vWarehouse.LocationID,
		ToWarehouseID:              vWarehouse.ID,
		ProductID:                  req.ProductID,
		StockmovementvehicleType:   model.StockmovementvehicleTypeIn,
		Notes:                      req.Notes,
		ReceivedGrossQuantity:      req.NetQuantity,
		ReceivedTareQuantity:       0,
		ReceivedNetQuantity:        req.NetQuantity,
		ReceivedTime:               &now,
		Shrinkage:                  0,
		StockmovementvehicleStatus: model.StockmovementvehicleStatusUnloading,
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

func (u usecase) Delete(loginUser jwt.UserLogin, id string) (err error) {
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

	if tStockmovementvehicle.StockmovementvehicleStatus == model.StockmovementvehicleStatusCompleted {
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

	if tStockmovementvehicle.StockmovementvehicleStatus == model.StockmovementvehicleStatusCompleted {
		return errors.New(fmt.Sprintf("unable to update data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tStockmovementvehicle.ToWarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tStockmovementvehicle.StockmovementvehicleStatus = model.StockmovementvehicleStatusCompleted
	err = u.stockmovementvehicleRepository.Save(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

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

func NewUsecase(stockmovementvehicleRepository stockmovementvehicle.Repository, warehouseRepository warehouse.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository) Usecase {
	return &usecase{
		stockmovementvehicleRepository: stockmovementvehicleRepository,
		warehouseRepository:            warehouseRepository,
		stockRepository:                stockRepository,
		stocklogRepository:             stocklogRepository,
	}
}
