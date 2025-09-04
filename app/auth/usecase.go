package auth

import (
	"errors"
	"time"

	"github.com/jihanlugas/warehouse/app/user"
	"github.com/jihanlugas/warehouse/app/userprovider"
	"github.com/jihanlugas/warehouse/app/warehouse"
	"github.com/jihanlugas/warehouse/config"
	"github.com/jihanlugas/warehouse/constant"
	"github.com/jihanlugas/warehouse/cryption"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
)

type Usecase interface {
	SignIn(req request.Signin) (token string, userLogin jwt.UserLogin, err error)
	RefreshToken(userLogin jwt.UserLogin) (token string, err error)
	Init(userLogin jwt.UserLogin) (vUser model.UserView, vWarehouse model.WarehouseView, err error)
	GoogleCallback(callback response.GoogleCallback) (token string, userLogin jwt.UserLogin, err error)
	GoogleLinkCallback(loginUser jwt.UserLogin, callback response.GoogleCallback) (err error)
	GoogleUnlink(loginUser jwt.UserLogin) (err error)
}

type usecase struct {
	userRepository         user.Repository
	warehouseRepository    warehouse.Repository
	userproviderRepository userprovider.Repository
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
	userLogin.LocationID = tWarehouse.LocationID
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

func (u usecase) GoogleCallback(callback response.GoogleCallback) (token string, userLogin jwt.UserLogin, err error) {
	var tUserprovider model.Userprovider
	var tWarehouse model.Warehouse

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUserprovider, err = u.userproviderRepository.GetTableByProviderNameAndEmail(conn, constant.OauthProviderGoogle, callback.Email, "User")
	if err != nil {
		err = errors.New("akun tidak terdaftar")
		return token, userLogin, err
	}

	if tUserprovider.User == nil {
		err = errors.New("akun tidak valid")
		return token, userLogin, err
	}

	if !tUserprovider.User.IsActive {
		return token, userLogin, errors.New("user not active")
	}

	if tUserprovider.User.UserRole == model.UserRoleOperator {
		tWarehouse, err = u.warehouseRepository.GetTableById(conn, tUserprovider.User.WarehouseID)
		if err != nil {
			return token, userLogin, errors.New("warehouse not found : " + err.Error())
		}
	}

	now := time.Now()
	tx := conn.Begin()

	tUserprovider.User.LastLoginDt = &now
	tUserprovider.User.UpdateBy = tUserprovider.User.ID
	err = u.userRepository.Update(tx, model.User{
		ID:          tUserprovider.User.ID,
		LastLoginDt: &now,
		UpdateBy:    tUserprovider.User.ID,
	})
	if err != nil {
		return token, userLogin, err
	}

	expiredAt := time.Now().Add(time.Minute * time.Duration(config.AuthTokenExpiredMinute))
	userLogin.ExpiredDt = expiredAt
	userLogin.UserID = tUserprovider.User.ID
	userLogin.UserRole = tUserprovider.User.UserRole
	userLogin.PassVersion = tUserprovider.User.PassVersion
	userLogin.LocationID = tWarehouse.LocationID
	userLogin.WarehouseID = tWarehouse.ID
	token, err = jwt.CreateToken(userLogin)
	if err != nil {
		return token, userLogin, err
	}

	err = tx.Commit().Error
	if err != nil {
		return token, userLogin, err
	}

	return token, userLogin, err
}

func (u usecase) GoogleLinkCallback(loginUser jwt.UserLogin, callback response.GoogleCallback) (err error) {
	var tUserprovider model.Userprovider
	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUserprovider, err = u.userproviderRepository.GetTableByProviderNameAndEmail(conn, constant.OauthProviderGoogle, callback.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if tUserprovider.ID != "" {
		return errors.New("email sudah terdaftar")
	}

	tx := conn.Begin()

	err = u.userproviderRepository.Create(tx, model.Userprovider{
		UserID:         loginUser.UserID,
		ProviderName:   constant.OauthProviderGoogle,
		ProviderUserID: callback.ID,
		Email:          callback.Email,
		Fullname:       callback.Name,
		CreateBy:       loginUser.UserID,
		UpdateBy:       loginUser.UserID,
	})
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err

}

func (u usecase) GoogleUnlink(loginUser jwt.UserLogin) (err error) {
	var tUserprovider model.Userprovider

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tUserprovider, err = u.userproviderRepository.GetTableByProviderNameAndUserId(conn, "google", loginUser.UserID)
	if err != nil {
		return err
	}

	tx := conn.Begin()

	err = u.userproviderRepository.Delete(tx, tUserprovider)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err

}

func NewUsecase(userRepository user.Repository, warehouseRepository warehouse.Repository, userproviderRepository userprovider.Repository) Usecase {
	return usecase{
		userRepository:         userRepository,
		warehouseRepository:    warehouseRepository,
		userproviderRepository: userproviderRepository,
	}
}
