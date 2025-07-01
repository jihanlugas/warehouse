package vehicle

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tVehicle model.Vehicle, err error)
	GetTableByName(conn *gorm.DB, name string) (tVehicle model.Vehicle, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vVehicle model.VehicleView, err error)
	GetViewByName(conn *gorm.DB, name string) (vVehicle model.VehicleView, err error)
	Create(conn *gorm.DB, tVehicle model.Vehicle) error
	Update(conn *gorm.DB, tVehicle model.Vehicle) error
	Save(conn *gorm.DB, tVehicle model.Vehicle) error
	Delete(conn *gorm.DB, tVehicle model.Vehicle) error
	Page(conn *gorm.DB, req request.PageVehicle) (vVehicles []model.VehicleView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "vehicle"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tVehicle model.Vehicle, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tVehicle).Error
	return tVehicle, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tVehicle model.Vehicle, err error) {
	err = conn.Where("name = ? ", name).First(&tVehicle).Error
	return tVehicle, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vVehicle model.VehicleView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vVehicle).Error
	return vVehicle, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vVehicle model.VehicleView, err error) {
	err = conn.Where("name = ? ", name).First(&vVehicle).Error
	return vVehicle, err
}

func (r repository) Create(conn *gorm.DB, tVehicle model.Vehicle) error {
	return conn.Create(&tVehicle).Error
}

func (r repository) Update(conn *gorm.DB, tVehicle model.Vehicle) error {
	return conn.Model(&tVehicle).Updates(&tVehicle).Error
}

func (r repository) Save(conn *gorm.DB, tVehicle model.Vehicle) error {
	return conn.Save(&tVehicle).Error
}

func (r repository) Delete(conn *gorm.DB, tVehicle model.Vehicle) error {
	return conn.Delete(&tVehicle).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageVehicle) (vVehicles []model.VehicleView, count int64, err error) {
	query := conn.Model(&vVehicles)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.PlateNumber != "" {
		query = query.Where("plate_number ILIKE ?", "%"+req.PlateNumber+"%")
	}
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.NIK != "" {
		query = query.Where("nik ILIKE ?", "%"+req.NIK+"%")
	}
	if req.DriverName != "" {
		query = query.Where("driver_name ILIKE ?", "%"+req.DriverName+"%")
	}
	if req.PhoneNumber != "" {
		query = query.Where("phone_number ILIKE ?", "%"+utils.FormatPhoneTo62(req.PhoneNumber)+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vVehicles, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vVehicles).Error
	if err != nil {
		return vVehicles, count, err
	}

	return vVehicles, count, err
}

func NewRepository() Repository {
	return repository{}
}
