package retail

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/customer"
	"github.com/jihanlugas/warehouse/app/retailproduct"
	"github.com/jihanlugas/warehouse/app/stock"
	"github.com/jihanlugas/warehouse/app/stocklog"
	"github.com/jihanlugas/warehouse/app/vehicle"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageRetail) (vRetails []model.RetailView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vRetail model.RetailView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateRetail) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateRetail) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	retailRepository        Repository
	retailproductRepository retailproduct.Repository
	stockRepository         stock.Repository
	stocklogRepository      stocklog.Repository
	customerRepository      customer.Repository
	warehouseRepository     warehouse.Repository
	vehicleRepository       vehicle.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageRetail) (vRetails []model.RetailView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vRetails, count, err = u.retailRepository.Page(conn, req)
	if err != nil {
		return vRetails, count, err
	}

	return vRetails, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vRetail model.RetailView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vRetail, err = u.retailRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vRetail, errors.New(fmt.Sprint("failed to get retail: ", err))
	}

	return vRetail, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateRetail) error {
	var err error
	var tCustomer model.Customer
	var tRetail model.Retail
	var tRetailproduct model.Retailproduct

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	if req.IsNewCustomer {
		tCustomer = model.Customer{
			ID:          utils.GetUniqueID(),
			Name:        req.CustomerName,
			PhoneNumber: utils.FormatPhoneTo62(req.CustomerPhoneNumber),
			CreateBy:    loginUser.UserID,
			UpdateBy:    loginUser.UserID,
		}

		err = u.customerRepository.Create(tx, tCustomer)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create %s: %v", u.customerRepository.Name(), err))
		}
		req.CustomerID = tCustomer.ID
	}

	tRetail = model.Retail{
		ID:         utils.GetUniqueID(),
		CustomerID: req.CustomerID,
		Notes:      req.Notes,
		CreateBy:   loginUser.UserID,
		UpdateBy:   loginUser.UserID,
	}

	err = u.retailRepository.Create(tx, tRetail)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create retail: ", err))
	}

	for _, product := range req.Products {
		tRetailproduct = model.Retailproduct{
			ID:        utils.GetUniqueID(),
			RetailID:  tRetail.ID,
			ProductID: product.ProductID,
			UnitPrice: product.UnitPrice,
			CreateBy:  loginUser.UserID,
			UpdateBy:  loginUser.UserID,
		}
		err = u.retailproductRepository.Create(tx, tRetailproduct)
		if err != nil {
			return fmt.Errorf("failed to create %s: %v", u.retailproductRepository.Name(), err)
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateRetail) error {
	var err error
	var tRetail model.Retail

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tRetail, err = u.retailRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get retail: ", err))
	}

	tx := conn.Begin()
	tRetail.Notes = req.Notes
	tRetail.UpdateBy = loginUser.UserID
	err = u.retailRepository.Save(tx, tRetail)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update retail: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tRetail model.Retail

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tRetail, err = u.retailRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprint("failed to get retail: ", err))
	}

	tx := conn.Begin()

	err = u.retailRepository.Delete(tx, tRetail)
	if err != nil {
		return errors.New(fmt.Sprint("failed to delete retail: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(retailRepository Repository, retailproductRepository retailproduct.Repository, stockRepository stock.Repository, stocklogRepository stocklog.Repository, customerRepository customer.Repository, warehouseRepository warehouse.Repository, vehicleRepository vehicle.Repository) Usecase {
	return &usecase{
		retailRepository:        retailRepository,
		retailproductRepository: retailproductRepository,
		stockRepository:         stockRepository,
		stocklogRepository:      stocklogRepository,
		customerRepository:      customerRepository,
		warehouseRepository:     warehouseRepository,
		vehicleRepository:       vehicleRepository,
	}
}
