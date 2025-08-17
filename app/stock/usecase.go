package stock

import (
	"errors"
	"fmt"

	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
)

type Usecase interface {
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStock model.StockView, err error)
	Update(loginUser jwt.UserLogin, id string, req request.UpdateStock) (err error)
}

type usecase struct {
	stockRepository    Repository
	stocklogRepository stocklog.Repository
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStock model.StockView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStock, err = u.stockRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vStock, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vStock.WarehouseID) {
		return vStock, errors.New(response.ErrorHandlerIDOR)
	}

	return vStock, err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateStock) (err error) {
	var tStock model.Stock
	var tStocklog model.Stocklog

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStock, err = u.stockRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockRepository.Name(), err))
	}

	tx := conn.Begin()
	tStock.Quantity = req.Quantity
	tStock.UpdateBy = loginUser.UserID
	err = u.stockRepository.Save(tx, tStock)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.stockRepository.Name(), err))
	}

	tStocklog = model.Stocklog{
		WarehouseID:     tStock.WarehouseID,
		StockID:         tStock.ID,
		ProductID:       tStock.ProductID,
		StocklogType:    model.StocklogTypeAdjustment,
		NetQuantity:     tStock.Quantity,
		CurrentQuantity: tStock.Quantity,
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

func NewUsecase(stockRepository Repository, stocklogRepository stocklog.Repository) Usecase {
	return &usecase{
		stockRepository:    stockRepository,
		stocklogRepository: stocklogRepository,
	}
}
