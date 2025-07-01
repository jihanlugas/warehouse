package product

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tProduct model.Product, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vProduct model.ProductView, err error)
	Create(conn *gorm.DB, tProduct model.Product) error
	Update(conn *gorm.DB, tProduct model.Product) error
	Save(conn *gorm.DB, tProduct model.Product) error
	Delete(conn *gorm.DB, tProduct model.Product) error
	Page(conn *gorm.DB, req request.PageProduct) (vProducts []model.ProductView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "retail"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tProduct model.Product, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tProduct).Error
	return tProduct, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tProduct model.Product, err error) {
	err = conn.Where("name = ? ", name).First(&tProduct).Error
	return tProduct, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vProduct model.ProductView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vProduct).Error
	return vProduct, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vProduct model.ProductView, err error) {
	err = conn.Where("name = ? ", name).First(&vProduct).Error
	return vProduct, err
}

func (r repository) Create(conn *gorm.DB, tProduct model.Product) error {
	return conn.Create(&tProduct).Error
}

func (r repository) Update(conn *gorm.DB, tProduct model.Product) error {
	return conn.Model(&tProduct).Updates(&tProduct).Error
}

func (r repository) Save(conn *gorm.DB, tProduct model.Product) error {
	return conn.Save(&tProduct).Error
}

func (r repository) Delete(conn *gorm.DB, tProduct model.Product) error {
	return conn.Delete(&tProduct).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageProduct) (vProducts []model.ProductView, count int64, err error) {
	query := conn.Model(&vProducts)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
	}
	if req.StartCreateDt != nil {
		query = query.Where("create_dt >= ?", req.StartCreateDt)
	}
	if req.EndCreateDt != nil {
		query = query.Where("create_dt <= ?", req.EndCreateDt)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vProducts, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vProducts).Error
	if err != nil {
		return vProducts, count, err
	}

	return vProducts, count, err
}

func NewRepository() Repository {
	return repository{}
}
