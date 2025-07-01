package stock

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tStock model.Stock, err error)
	GetTableByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (tStock model.Stock, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vStock model.StockView, err error)
	GetViewByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (vStock model.StockView, err error)
	Create(conn *gorm.DB, tStock model.Stock) error
	Update(conn *gorm.DB, tStock model.Stock) error
	Save(conn *gorm.DB, tStock model.Stock) error
	Delete(conn *gorm.DB, tStock model.Stock) error
	Page(conn *gorm.DB, req request.PageStock) (vStocks []model.StockView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "stock"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tStock model.Stock, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tStock).Error
	return tStock, err
}

func (r repository) GetTableByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (tStock model.Stock, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("warehouse_id = ? ", warehouseID).Where("product_id = ? ", productID).First(&tStock).Error
	return tStock, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vStock model.StockView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vStock).Error
	return vStock, err
}

func (r repository) GetViewByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (vStock model.StockView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("warehouse_id = ? ", warehouseID).Where("product_id = ? ", productID).First(&vStock).Error
	return vStock, err
}

func (r repository) Create(conn *gorm.DB, tStock model.Stock) error {
	return conn.Create(&tStock).Error
}

func (r repository) Update(conn *gorm.DB, tStock model.Stock) error {
	return conn.Model(&tStock).Updates(&tStock).Error
}

func (r repository) Save(conn *gorm.DB, tStock model.Stock) error {
	return conn.Save(&tStock).Error
}

func (r repository) Delete(conn *gorm.DB, tStock model.Stock) error {
	return conn.Delete(&tStock).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageStock) (vStocks []model.StockView, count int64, err error) {
	query := conn.Model(&vStocks)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.StartQuantity != nil {
		query = query.Where("quantity >= ?", req.StartQuantity)
	}
	if req.EndQuantity != nil {
		query = query.Where("quantity <= ?", req.EndQuantity)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vStocks, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vStocks).Error
	if err != nil {
		return vStocks, count, err
	}

	return vStocks, count, err
}

func NewRepository() Repository {
	return repository{}
}
