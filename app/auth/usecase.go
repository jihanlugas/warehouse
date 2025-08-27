package auth

import (
	"errors"
	"time"

	"github.com/jihanlugas/warehouse/app/user"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/config"
	"github.com/jihanlugas/warehouse/cryption"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
)

type Usecase interface {
	SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error)
	RefreshToken(userLogin jwt.UserLogin) (token string, err error)
	Init(userLogin jwt.UserLogin) (vUser model.UserView, vWarehouse model.WarehouseView, err error)
}

type usecase struct {
	userRepository      user.Repository
	warehouseRepository warehouse.Repository
}

func (u usecase) SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error) {
	var tUser model.User
	var tWarehouse model.Warehouse

	conn, closeConn := db.GetConnection()
	defer closeConn()

	if utils.IsValidEmail(req.Username) {
		tUser, err = u.userRepository.GetByEmail(conn, req.Username)
	} else {
		tUser, err = u.userRepository.GetByUsername(conn, req.Username)
	}

	if err != nil {
		return "", userLogin, err
	}

	err = cryption.CheckAES64(req.Passwd, tUser.Passwd)
	if err != nil {
		return "", userLogin, errors.New("invalid username or password")
	}

	if !tUser.IsActive {
		return "", userLogin, errors.New("user not active")
	}

	if tUser.UserRole == model.UserRoleOperator {
		tWarehouse, err = u.warehouseRepository.GetTableById(conn, tUser.WarehouseID)
		if err != nil {
			return "", userLogin, errors.New("warehouse not found : " + err.Error())
		}
	}

	now := time.Now()
	tx := conn.Begin()

	tUser.LastLoginDt = &now
	tUser.UpdateBy = tUser.ID
	err = u.userRepository.Update(tx, model.User{
		ID:          tUser.ID,
		LastLoginDt: &now,
		UpdateBy:    tUser.ID,
	})
	if err != nil {
		return "", userLogin, err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", userLogin, err
	}

	expiredAt := time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))
	userLogin.ExpiredDt = expiredAt
	userLogin.UserID = tUser.ID
	userLogin.UserRole = tUser.UserRole
	userLogin.PassVersion = tUser.PassVersion
	userLogin.WarehouseID = tWarehouse.ID
	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return "", userLogin, err
	}

	return token, userLogin, err
}

func (u usecase) RefreshToken(userLogin jwt.UserLogin) (token string, err error) {
	userLogin.ExpiredDt = time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))

	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return token, err
	}

	return token, err
}

func (u usecase) Init(userLogin jwt.UserLogin) (vUser model.UserView, vWarehouse model.WarehouseView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vUser, err = u.userRepository.GetViewById(conn, userLogin.UserID, "Userprivilege", "Userproviders")
	if err != nil {
		return vUser, vWarehouse, err
	}

	if vUser.UserRole == model.UserRoleOperator {
		vWarehouse, err = u.warehouseRepository.GetViewById(conn, userLogin.WarehouseID)
		if err != nil {
			return vUser, vWarehouse, err
		}
	}

	return vUser, vWarehouse, err
}

func NewUsecase(userRepository user.Repository, warehouseRepository warehouse.Repository) Usecase {
	return usecase{
		userRepository:      userRepository,
		warehouseRepository: warehouseRepository,
	}
}
