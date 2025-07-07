package product

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
	Page(loginUser jwt.UserLogin, req request.PageProduct) (vProducts []model.ProductView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vProduct model.ProductView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateProduct) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateProduct) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	productRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageProduct) (vProducts []model.ProductView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vProducts, count, err = u.productRepository.Page(conn, req)
	if err != nil {
		return vProducts, count, err
	}

	return vProducts, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vProduct model.ProductView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vProduct, err = u.productRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vProduct, errors.New(fmt.Sprintf("failed to get %s: %v", u.productRepository.Name(), err))
	}

	return vProduct, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateProduct) error {
	var err error
	var tProduct model.Product

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tProduct = model.Product{
		ID:          utils.GetUniqueID(),
		Name:        req.Name,
		Description: req.Description,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.productRepository.Create(tx, tProduct)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.productRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateProduct) error {
	var err error
	var tProduct model.Product

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tProduct, err = u.productRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.productRepository.Name(), err))
	}

	tx := conn.Begin()
	tProduct.Name = req.Name
	tProduct.Description = req.Description
	tProduct.UpdateBy = loginUser.UserID
	err = u.productRepository.Save(tx, tProduct)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.productRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tProduct model.Product

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tProduct, err = u.productRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.productRepository.Name(), err))
	}

	tx := conn.Begin()

	err = u.productRepository.Delete(tx, tProduct)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.productRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(productRepository Repository) Usecase {
	return &usecase{
		productRepository: productRepository,
	}
}
