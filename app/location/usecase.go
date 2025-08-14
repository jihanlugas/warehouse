package location

import (
	"errors"
	"fmt"

	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageLocation) (vLocations []model.LocationView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vLocation model.LocationView, err error)
}

type usecase struct {
	locationRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageLocation) (vLocations []model.LocationView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vLocations, count, err = u.locationRepository.Page(conn, req)
	if err != nil {
		return vLocations, count, err
	}

	return vLocations, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vLocation model.LocationView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vLocation, err = u.locationRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vLocation, errors.New(fmt.Sprintf("failed to get %s: %v", u.locationRepository.Name(), err))
	}

	return vLocation, err
}

func NewUsecase(locationRepository Repository) Usecase {
	return &usecase{
		locationRepository: locationRepository,
	}
}
