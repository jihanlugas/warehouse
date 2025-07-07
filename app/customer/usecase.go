package customer

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
	Page(loginUser jwt.UserLogin, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vCustomer model.CustomerView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateCustomer) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateCustomer) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	customerRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vCustomers, count, err = u.customerRepository.Page(conn, req)
	if err != nil {
		return vCustomers, count, err
	}

	return vCustomers, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vCustomer model.CustomerView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vCustomer, err = u.customerRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vCustomer, errors.New(fmt.Sprintf("failed to get %s: %v", u.customerRepository.Name(), err))
	}

	return vCustomer, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateCustomer) error {
	var err error
	var tCustomer model.Customer

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tCustomer = model.Customer{
		ID:          utils.GetUniqueID(),
		Name:        req.Name,
		PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
		Email:       req.Email,
		Address:     req.Address,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.customerRepository.Create(tx, tCustomer)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.customerRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateCustomer) error {
	var err error
	var tCustomer model.Customer

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tCustomer, err = u.customerRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.customerRepository.Name(), err))
	}

	tx := conn.Begin()

	tCustomer.Name = req.Name
	tCustomer.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tCustomer.Email = req.Email
	tCustomer.Address = req.Address
	tCustomer.UpdateBy = loginUser.UserID
	err = u.customerRepository.Save(tx, tCustomer)
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
	var tCustomer model.Customer

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tCustomer, err = u.customerRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.customerRepository.Name(), err))
	}

	tx := conn.Begin()

	err = u.customerRepository.Delete(tx, tCustomer)
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
