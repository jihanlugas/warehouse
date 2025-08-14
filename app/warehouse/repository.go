package warehouse

import (
	"fmt"
	"strings"

	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tWarehouse model.Warehouse, err error)
	GetTableByName(conn *gorm.DB, name string) (tWarehouse model.Warehouse, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vWarehouse model.WarehouseView, err error)
	GetViewByName(conn *gorm.DB, name string) (vWarehouse model.WarehouseView, err error)
	Create(conn *gorm.DB, tWarehouse model.Warehouse) error
	Update(conn *gorm.DB, tWarehouse model.Warehouse) error
	Save(conn *gorm.DB, tWarehouse model.Warehouse) error
	Delete(conn *gorm.DB, tWarehouse model.Warehouse) error
	Page(conn *gorm.DB, req request.PageWarehouse) (vWarehouses []model.WarehouseView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "warehouse"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tWarehouse model.Warehouse, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tWarehouse).Error
	return tWarehouse, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tWarehouse model.Warehouse, err error) {
	err = conn.Where("name = ? ", name).First(&tWarehouse).Error
	return tWarehouse, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vWarehouse model.WarehouseView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vWarehouse).Error
	return vWarehouse, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vWarehouse model.WarehouseView, err error) {
	err = conn.Where("name = ? ", name).First(&vWarehouse).Error
	return vWarehouse, err
}

func (r repository) Create(conn *gorm.DB, tWarehouse model.Warehouse) error {
	return conn.Create(&tWarehouse).Error
}

func (r repository) Update(conn *gorm.DB, tWarehouse model.Warehouse) error {
	return conn.Model(&tWarehouse).Updates(&tWarehouse).Error
}

func (r repository) Save(conn *gorm.DB, tWarehouse model.Warehouse) error {
	return conn.Save(&tWarehouse).Error
}

func (r repository) Delete(conn *gorm.DB, tWarehouse model.Warehouse) error {
	return conn.Delete(&tWarehouse).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageWarehouse) (vWarehouses []model.WarehouseView, count int64, err error) {
	query := conn.Model(&vWarehouses)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.LocationID != "" {
		query = query.Where("location_id = ?", req.LocationID)
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
		return vWarehouses, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "name", "asc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vWarehouses).Error
	if err != nil {
		return vWarehouses, count, err
	}

	return vWarehouses, count, err
}

func NewRepository() Repository {
	return repository{}
}
