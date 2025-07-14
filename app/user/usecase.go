package user

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/app/userprivilege"
	"github.com/jihanlugas/warehouse/cryption"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"time"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageUser) (vUsers []model.UserView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vUser model.UserView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateUser) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateUser) error
	UpdateUserprivilege(loginUser jwt.UserLogin, id string, req request.UpdateUserprivilege) error
	ChangePassword(loginUser jwt.UserLogin, req request.ChangePassword) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	userRepository          Repository
	userprivilegeRepository userprivilege.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageUser) (vUsers []model.UserView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	if jwt.IsSaveWarehouseIDOR(loginUser, req.WarehouseID) {
		return vUsers, count, errors.New(response.ErrorHandlerIDOR)
	}

	vUsers, count, err = u.userRepository.Page(conn, req)
	if err != nil {
		return vUsers, count, err
	}

	return vUsers, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vUser model.UserView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vUser, err = u.userRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vUser, errors.New(fmt.Sprintf("failed to get %s: %v", u.userRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vUser.WarehouseID) {
		return vUser, errors.New(response.ErrorHandlerIDOR)
	}

	return vUser, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateUser) error {
	var err error
	var tUser model.User
	var tUserprivilege model.Userprivilege

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	now := time.Now()

	encodePasswd, err := cryption.EncryptAES64(req.Passwd)
	if err != nil {
		return errors.New(fmt.Sprint("failed to encode password: ", err))
	}

	tUser = model.User{
		ID:                utils.GetUniqueID(),
		WarehouseID:       req.WarehouseID,
		Role:              model.UserRoleOperator,
		Email:             req.Email,
		Username:          req.Username,
		PhoneNumber:       utils.FormatPhoneTo62(req.PhoneNumber),
		Address:           req.Address,
		Fullname:          req.Fullname,
		Passwd:            encodePasswd,
		PassVersion:       1,
		IsActive:          true,
		PhotoID:           "",
		LastLoginDt:       nil,
		BirthDt:           req.BirthDt,
		BirthPlace:        req.BirthPlace,
		AccountVerifiedDt: &now,
		CreateBy:          loginUser.UserID,
		UpdateBy:          loginUser.UserID,
	}

	err = u.userRepository.Create(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.userRepository.Name(), err))
	}

	tUserprivilege = model.Userprivilege{
		UserID:        tUser.ID,
		StockIn:       req.StockIn,
		TransferOut:   req.TransferOut,
		TransferIn:    req.TransferIn,
		PurchaseOrder: req.PurchaseOrder,
		Retail:        req.Retail,
		CreateBy:      loginUser.UserID,
		UpdateBy:      loginUser.UserID,
	}
	err = u.userprivilegeRepository.Create(tx, tUserprivilege)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.userprivilegeRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateUser) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.userRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tUser.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	if tUser.Email != req.Email {
		_, err = u.userRepository.GetByEmail(tx, req.Email)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to update %s: %v", u.userRepository.Name(), err)
			}
		} else {
			return errors.New("email already exist")
		}
	}

	if tUser.PhoneNumber != utils.FormatPhoneTo62(req.PhoneNumber) {
		_, err = u.userRepository.GetByPhoneNumber(tx, utils.FormatPhoneTo62(req.PhoneNumber))
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to update %s: %v", u.userRepository.Name(), err)
			}
		} else {
			return errors.New("phone number already exist")
		}
	}

	tUser.Fullname = req.Fullname
	tUser.Email = req.Email
	tUser.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tUser.Username = req.Username
	tUser.Address = req.Address
	tUser.BirthDt = req.BirthDt
	tUser.BirthPlace = req.BirthPlace
	tUser.UpdateBy = loginUser.UserID
	err = u.userRepository.Save(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.userRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) UpdateUserprivilege(loginUser jwt.UserLogin, id string, req request.UpdateUserprivilege) error {
	var err error
	var tUser model.User
	var tUserprivilege model.Userprivilege

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, id, "Userprivilege")
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.userRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tUser.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	if tUser.Userprivilege.ID == "" {
		return errors.New(response.ErrorDataNotFound)
	}

	tx := conn.Begin()

	tUserprivilege = *tUser.Userprivilege
	tUserprivilege.StockIn = req.StockIn
	tUserprivilege.TransferOut = req.TransferOut
	tUserprivilege.TransferIn = req.TransferIn
	tUserprivilege.PurchaseOrder = req.PurchaseOrder
	tUserprivilege.Retail = req.Retail
	err = u.userprivilegeRepository.Save(conn, tUserprivilege)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update %s: %v", u.userprivilegeRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) ChangePassword(loginUser jwt.UserLogin, req request.ChangePassword) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, loginUser.UserID)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.userRepository.Name(), err))
	}

	tx := conn.Begin()

	err = cryption.CheckAES64(req.CurrentPasswd, tUser.Passwd)
	if err != nil {
		return errors.New(fmt.Sprint("invalid current password"))
	}

	encodePasswd, err := cryption.EncryptAES64(req.Passwd)
	if err != nil {
		return errors.New(fmt.Sprint("failed to encode password: ", err))
	}

	tUser.Passwd = encodePasswd
	tUser.PassVersion += 1
	tUser.UpdateBy = loginUser.UserID
	err = u.userRepository.Save(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprint("failed to update password: ", err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tUser model.User

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUser, err = u.userRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.userRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tUser.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.userRepository.Delete(tx, tUser)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.userRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(userRepository Repository, userprivilegeRepository userprivilege.Repository) Usecase {
	return &usecase{
		userRepository:          userRepository,
		userprivilegeRepository: userprivilegeRepository,
	}
}
