package stockin

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/stockmovement"
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
	Page(loginUser jwt.UserLogin, req request.PageStockin) (vStockins []model.StockinView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockin model.StockinView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateStockin) error
}

type usecase struct {
	stockinRepository       Repository
	warehouseRepository     warehouse.Repository
	stockRepository         stock.Repository
	stocklogRepository      stocklog.Repository
	stockmovementRepository stockmovement.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageStockin) (vStockins []model.StockinView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockins, count, err = u.stockinRepository.Page(conn, req)
	if err != nil {
		return vStockins, count, err
	}

	return vStockins, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockin model.StockinView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockin, err = u.stockinRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vStockin, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockinRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStockin.WarehouseID) {
		return vStockin, errors.New(response.ErrorHandlerIDOR)
	}

	return vStockin, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateStockin) error {
	var err error
	var vWarehouse model.WarehouseView
	var tStock model.Stock
	var tStockmovement model.Stockmovement
	var tStocklog model.Stocklog

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	vWarehouse, err = u.warehouseRepository.GetViewById(conn, req.WarehouseID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.warehouseRepository.Name(), err))
	}

	if !vWarehouse.IsStockin {
		return errors.New(fmt.Sprint("this warehouse is not allowed to create stockin"))
	}

	tStock, err = u.stockRepository.GetTableByWarehouseIdAndProductId(tx, req.WarehouseID, req.ProductID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockRepository.Name(), err))
		}
		tStock = model.Stock{
			ID:          utils.GetUniqueID(),
			WarehouseID: req.WarehouseID,
			ProductID:   req.ProductID,
			Quantity:    0,
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}
		err = u.stockRepository.Save(tx, tStock)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockRepository.Name(), err))
		}
	}

	tStockmovement = model.Stockmovement{
		ID:            utils.GetUniqueID(),
		ToWarehouseID: req.WarehouseID,
		ProductID:     req.ProductID,
		Type:          model.StockMovementTypeIn,
		Remark:        req.Remark,
		CreateBy:      loginUser.UserID,
		UpdateBy:      loginUser.UserID,
	}
	err = u.stockmovementRepository.Create(tx, tStockmovement)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockmovementRepository.Name(), err))
	}

	CurrentQuantity := 0.0
	CurrentQuantity = tStock.Quantity + req.NetQuantity
	tStock.Quantity = CurrentQuantity
	tStock.UpdateBy = loginUser.UserID
	err = u.stockRepository.Save(tx, tStock)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.stockRepository.Name(), err))
	}

	tStocklog = model.Stocklog{
		WarehouseID:     req.WarehouseID,
		StockID:         tStock.ID,
		StockmovementID: tStockmovement.ID,
		ProductID:       req.ProductID,
		Type:            model.StockLogTypeIn,
		GrossQuantity:   req.GrossQuantity,
		TareQuantity:    req.TareQuantity,
		NetQuantity:     req.NetQuantity,
		CurrentQuantity: CurrentQuantity,
		CreateBy:        loginUser.UserID,
		UpdateBy:        loginUser.UserID,
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

func NewUsecase(stockinRepository Repository, warehouseRepository warehouse.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, stockmovementRepository stockmovement.Repository) Usecase {
	return &usecase{
		stockinRepository:       stockinRepository,
		warehouseRepository:     warehouseRepository,
		stockRepository:         stockRepository,
		stocklogRepository:      stocklogRepository,
		stockmovementRepository: stockmovementRepository,
	}
}
