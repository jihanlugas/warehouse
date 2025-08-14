package stockmovementvehicle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jihanlugas/warehouse/app/photo"
	"github.com/jihanlugas/warehouse/app/stockmovementvehiclephoto"
	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	Delete(loginUser jwt.UserLogin, id string) (err error)
	UploadPhoto(loginUser jwt.UserLogin, id string, req request.CreateStockmovementvehiclephoto) (err error)
}

type usecase struct {
	stockmovementvehicleRepository      Repository
	stockmovementvehiclephotoRepository stockmovementvehiclephoto.Repository
	photoRepository                     photo.Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicles, count, err = u.stockmovementvehicleRepository.Page(conn, req)
	if err != nil {
		return vStockmovementvehicles, count, err
	}

	return vStockmovementvehicles, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vStockmovementvehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	//if jwt.IsSaveWarehouseIDOR(loginUser, vStockmovementvehicle.FromWarehouseID) {
	//	return vStockmovementvehicle, errors.New(response.ErrorHandlerIDOR)
	//}

	return vStockmovementvehicle, err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) (err error) {
	var tStockmovementvehicle model.Stockmovementvehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tStockmovementvehicle, err = u.stockmovementvehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	if tStockmovementvehicle.StockmovementvehicleStatus == model.StockmovementvehicleStatusCompleted {
		return errors.New(fmt.Sprintf("unable to delete data with status %s", strings.ToLower(string(tStockmovementvehicle.StockmovementvehicleStatus))))
	}

	tx := conn.Begin()

	err = u.stockmovementvehicleRepository.Delete(tx, tStockmovementvehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) UploadPhoto(loginUser jwt.UserLogin, id string, req request.CreateStockmovementvehiclephoto) (err error) {
	var vStockmovementvehicle model.StockmovementvehicleView

	conn, closeConn := db.GetConnection()
	defer closeConn()

	vStockmovementvehicle, err = u.stockmovementvehicleRepository.GetViewById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.stockmovementvehicleRepository.Name(), err))
	}

	switch vStockmovementvehicle.StockmovementvehicleType {
	case model.StockmovementvehicleTypeTransfer:
		switch vStockmovementvehicle.StockmovementvehicleStatus {
		case "LOADING":
			if req.WarehouseID != vStockmovementvehicle.FromWarehouseID {
				return errors.New("unable to upload photo")
			}
			break
		case "UNLOADING":
			if req.WarehouseID != vStockmovementvehicle.ToWarehouseID {
				return errors.New("unable to upload photo")
			}
			break
		default:
			return errors.New("unable to upload photo")
		}
		break
	case model.StockmovementvehicleTypePurchaseorder:
		if vStockmovementvehicle.StockmovementvehicleStatus != "LOADING" && vStockmovementvehicle.FromWarehouseID != req.WarehouseID {
			return errors.New("unable to upload photo")
		}
		break
	case model.StockmovementvehicleTypeRetail:
		if vStockmovementvehicle.StockmovementvehicleStatus != "LOADING" && vStockmovementvehicle.FromWarehouseID != req.WarehouseID {
			return errors.New("unable to upload photo")
		}
		break
	default:
		return errors.New("invalid stock movement type")
	}

	tx := conn.Begin()

	tPhoto, err := u.photoRepository.Upload(tx, req.Photo, model.PhotoRefStockmovementvehiclephoto)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload %s: %v", u.photoRepository.Name(), err))
	}

	tStockmovementvehiclephoto := model.Stockmovementvehiclephoto{
		ID:                     utils.GetUniqueID(),
		WarehouseID:            req.WarehouseID,
		StockmovementvehicleID: req.StockmovementvehicleID,
		PhotoID:                tPhoto.ID,
		CreateBy:               loginUser.UserID,
		UpdateBy:               loginUser.UserID,
	}
	err = u.stockmovementvehiclephotoRepository.Create(tx, tStockmovementvehiclephoto)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.stockmovementvehiclephotoRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(stockmovementvehicleRepository Repository, stockmovementvehiclephotoRepository stockmovementvehiclephoto.Repository, photoRepository photo.Repository) Usecase {
	return &usecase{
		stockmovementvehicleRepository:      stockmovementvehicleRepository,
		stockmovementvehiclephotoRepository: stockmovementvehiclephotoRepository,
		photoRepository:                     photoRepository,
	}
}
