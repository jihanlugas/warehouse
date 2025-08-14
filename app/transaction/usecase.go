package transaction

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
	Page(loginUser jwt.UserLogin, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vTransaction model.TransactionView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateTransaction) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateTransaction) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	transactionRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vTransactions, count, err = u.transactionRepository.Page(conn, req)
	if err != nil {
		return vTransactions, count, err
	}

	return vTransactions, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vTransaction model.TransactionView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vTransaction, err = u.transactionRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vTransaction, errors.New(fmt.Sprintf("failed to get %s: %v", u.transactionRepository.Name(), err))
	}

	return vTransaction, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateTransaction) error {
	var err error
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tTransaction = model.Transaction{
		ID:                 utils.GetUniqueID(),
		RelatedID:          req.RelatedID,
		TransactionRelated: model.TransactionRelated(req.TransactionRelated),
		TransactionType:    model.TransactionTypePayment,
		Amount:             req.Amount,
		Notes:              req.Notes,
		CreateBy:           loginUser.UserID,
		UpdateBy:           loginUser.UserID,
	}

	err = u.transactionRepository.Create(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.transactionRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateTransaction) error {
	var err error
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tTransaction, err = u.transactionRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.transactionRepository.Name(), err))
	}

	tx := conn.Begin()

	tTransaction.Notes = req.Notes
	tTransaction.UpdateBy = loginUser.UserID
	err = u.transactionRepository.Save(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.transactionRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tTransaction model.Transaction

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tTransaction, err = u.transactionRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.transactionRepository.Name(), err))
	}

	tx := conn.Begin()

	err = u.transactionRepository.Delete(tx, tTransaction)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.transactionRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(transactionRepository Repository) Usecase {
	return &usecase{
		transactionRepository: transactionRepository,
	}
}
