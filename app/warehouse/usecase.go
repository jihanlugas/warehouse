package warehouse

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageWarehouse) (vWarehouses []model.WarehouseView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vWarehouse model.WarehouseView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateWarehouse) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateWarehouse) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	customerRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageWarehouse) (vWarehouses []model.WarehouseView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vWarehouses, count, err = u.customerRepository.Page(conn, req)
	if err != nil {
		return vWarehouses, count, err
	}

	return vWarehouses, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vWarehouse model.WarehouseView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vWarehouse, err = u.customerRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vWarehouse, errors.New(fmt.Sprintf("failed to get %s: %v", u.customerRepository.Name(), err))
	}

	return vWarehouse, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateWarehouse) error {
	var err error
	var tWarehouse model.Warehouse

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tWarehouse = model.Warehouse{
		ID:              utils.GetUniqueID(),
		Name:            req.Name,
		Location:        req.Location,
		IsStockin:       false,
		IsInbound:       false,
		IsOutbound:      false,
		IsRetail:        false,
		IsPurchaseorder: false,
		CreateBy:        loginUser.UserID,
		UpdateBy:        loginUser.UserID,
	}

	err = u.customerRepository.Create(tx, tWarehouse)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.customerRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateWarehouse) error {
	var err error
	var tWarehouse model.Warehouse

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tWarehouse, err = u.customerRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.customerRepository.Name(), err))
	}

	tx := conn.Begin()

	tWarehouse.Name = req.Name
	tWarehouse.Location = req.Location
	tWarehouse.UpdateBy = loginUser.UserID
	err = u.customerRepository.Save(tx, tWarehouse)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.customerRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tWarehouse model.Warehouse

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tWarehouse, err = u.customerRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.customerRepository.Name(), err))
	}

	tx := conn.Begin()

	err = u.customerRepository.Delete(tx, tWarehouse)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.customerRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(customerRepository Repository) Usecase {
	return &usecase{
		customerRepository: customerRepository,
	}
}
