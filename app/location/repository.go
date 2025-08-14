package location

import (
	"fmt"
	"strings"

	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tLocation model.Location, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vLocation model.LocationView, err error)
	Page(conn *gorm.DB, req request.PageLocation) (vLocations []model.LocationView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "location"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tLocation model.Location, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tLocation).Error
	return tLocation, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vLocation model.LocationView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vLocation).Error
	return vLocation, err
}

func (r repository) Page(conn *gorm.DB, req request.PageLocation) (vLocations []model.LocationView, count int64, err error) {
	query := conn.Model(&vLocations)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Address != "" {
		query = query.Where("address ILIKE ?", "%"+req.Address+"%")
	}
	if req.Notes != "" {
		query = query.Where("notes ILIKE ?", "%"+req.Notes+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vLocations, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vLocations).Error
	if err != nil {
		return vLocations, count, err
	}

	return vLocations, count, err
}

func NewRepository() Repository {
	return repository{}
}
