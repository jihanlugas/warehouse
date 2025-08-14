package vehicle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/response"
	"github.com/jihanlugas/warehouse/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageVehicle) (vVehicles []model.VehicleView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vVehicle model.VehicleView, err error)
	Create(loginUser jwt.UserLogin, req request.CreateVehicle) error
	Update(loginUser jwt.UserLogin, id string, req request.UpdateVehicle) error
	Delete(loginUser jwt.UserLogin, id string) error
}

type usecase struct {
	vehicleRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageVehicle) (vVehicles []model.VehicleView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vVehicles, count, err = u.vehicleRepository.Page(conn, req)
	if err != nil {
		return vVehicles, count, err
	}

	return vVehicles, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vVehicle model.VehicleView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vVehicle, err = u.vehicleRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vVehicle, errors.New(fmt.Sprintf("failed to get %s: %v", u.vehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, vVehicle.WarehouseID) {
		return vVehicle, errors.New(response.ErrorHandlerIDOR)
	}

	return vVehicle, err
}

func (u usecase) Create(loginUser jwt.UserLogin, req request.CreateVehicle) error {
	var err error
	var tVehicle model.Vehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	tVehicle = model.Vehicle{
		ID:          utils.GetUniqueID(),
		WarehouseID: req.WarehouseID,
		PlateNumber: strings.ToUpper(req.PlateNumber),
		Name:        req.Name,
		NIK:         req.NIK,
		DriverName:  req.DriverName,
		PhoneNumber: utils.FormatPhoneTo62(req.PhoneNumber),
		Notes:       req.Notes,
		CreateBy:    loginUser.UserID,
		UpdateBy:    loginUser.UserID,
	}

	err = u.vehicleRepository.Create(tx, tVehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.vehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Update(loginUser jwt.UserLogin, id string, req request.UpdateVehicle) error {
	var err error
	var tVehicle model.Vehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tVehicle, err = u.vehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.vehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tVehicle.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	tVehicle.Name = req.Name
	tVehicle.Notes = req.Notes
	tVehicle.PlateNumber = strings.ToUpper(req.PlateNumber)
	tVehicle.NIK = req.NIK
	tVehicle.DriverName = req.DriverName
	tVehicle.PhoneNumber = utils.FormatPhoneTo62(req.PhoneNumber)
	tVehicle.UpdateBy = loginUser.UserID
	err = u.vehicleRepository.Save(tx, tVehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save %s: %v", u.vehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func (u usecase) Delete(loginUser jwt.UserLogin, id string) error {
	var err error
	var tVehicle model.Vehicle

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tVehicle, err = u.vehicleRepository.GetTableById(conn, id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get %s: %v", u.vehicleRepository.Name(), err))
	}

	if jwt.IsSaveWarehouseIDOR(loginUser, tVehicle.WarehouseID) {
		return errors.New(response.ErrorHandlerIDOR)
	}

	tx := conn.Begin()

	err = u.vehicleRepository.Delete(tx, tVehicle)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete %s: %v", u.vehicleRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(vehicleRepository Repository) Usecase {
	return &usecase{
		vehicleRepository: vehicleRepository,
	}
}
