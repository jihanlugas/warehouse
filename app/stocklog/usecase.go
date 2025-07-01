package stocklog

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageStocklog) (vStocklogs []model.StocklogView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStocklog model.StocklogView, err error)
}

type usecase struct {
	repository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageStocklog) (vStocklogs []model.StocklogView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStocklogs, count, err = u.repository.Page(conn, req)
	if err != nil {
		return vStocklogs, count, err
	}

	return vStocklogs, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStocklog model.StocklogView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStocklog, err = u.repository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vStocklog, errors.New(fmt.Sprintf("failed to get %s: %v", u.repository.Name(), err))
	}

	return vStocklog, err
}

func NewUsecase(repository Repository) Usecase {
	return &usecase{
		repository: repository,
	}
}
