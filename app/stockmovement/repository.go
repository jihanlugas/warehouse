package stockmovement

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tStockmovement model.Stockmovement, err error)
	GetTableByName(conn *gorm.DB, name string) (tStockmovement model.Stockmovement, err error)
	GetTableByFromWarehouseIDAndProductID(conn *gorm.DB, fromWarehouseId, productId string, preloads ...string) (tStockmovement model.Stockmovement, err error)
	GetTableByRelatedIDAndProductID(conn *gorm.DB, relatedId, productId string, preloads ...string) (tStockmovement model.Stockmovement, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockmovement model.StockmovementView, err error)
	GetViewByName(conn *gorm.DB, name string) (vStockmovement model.StockmovementView, err error)
	GetViewByFromWarehouseIDAndProductID(conn *gorm.DB, fromWarehouseId, productId string, preloads ...string) (vStockmovement model.StockmovementView, err error)
	GetViewByRelatedIDAndProductID(conn *gorm.DB, relatedId, productId string, preloads ...string) (vStockmovement model.StockmovementView, err error)
	Create(conn *gorm.DB, tStockmovement model.Stockmovement) error
	Update(conn *gorm.DB, tStockmovement model.Stockmovement) error
	Save(conn *gorm.DB, tStockmovement model.Stockmovement) error
	Delete(conn *gorm.DB, tStockmovement model.Stockmovement) error
	Page(conn *gorm.DB, req request.PageStockmovement) (vStockmovements []model.StockmovementView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "stockmovement"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tStockmovement model.Stockmovement, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tStockmovement).Error
	return tStockmovement, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tStockmovement model.Stockmovement, err error) {
	err = conn.Where("name = ? ", name).First(&tStockmovement).Error
	return tStockmovement, err
}

func (r repository) GetTableByFromWarehouseIDAndProductID(conn *gorm.DB, fromWarehouseId, productId string, preloads ...string) (tStockmovement model.Stockmovement, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("from_warehouse_id = ? ", fromWarehouseId).Where("product_id = ? ", productId).First(&tStockmovement).Error
	return tStockmovement, err
}

func (r repository) GetTableByRelatedIDAndProductID(conn *gorm.DB, relatedId, productId string, preloads ...string) (tStockmovement model.Stockmovement, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("related_id = ? ", relatedId).Where("product_id = ? ", productId).First(&tStockmovement).Error
	return tStockmovement, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockmovement model.StockmovementView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vStockmovement).Error
	return vStockmovement, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vStockmovement model.StockmovementView, err error) {
	err = conn.Where("name = ? ", name).First(&vStockmovement).Error
	return vStockmovement, err
}

func (r repository) GetViewByFromWarehouseIDAndProductID(conn *gorm.DB, fromWarehouseId, productId string, preloads ...string) (vStockmovement model.StockmovementView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("from_warehouse_id = ? ", fromWarehouseId).Where("product_id = ? ", productId).First(&vStockmovement).Error
	return vStockmovement, err
}

func (r repository) GetViewByRelatedIDAndProductID(conn *gorm.DB, relatedId, productId string, preloads ...string) (vStockmovement model.StockmovementView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("related_id = ? ", relatedId).Where("product_id = ? ", productId).First(&vStockmovement).Error
	return vStockmovement, err
}

func (r repository) Create(conn *gorm.DB, tStockmovement model.Stockmovement) error {
	return conn.Create(&tStockmovement).Error
}

func (r repository) Update(conn *gorm.DB, tStockmovement model.Stockmovement) error {
	return conn.Model(&tStockmovement).Updates(&tStockmovement).Error
}

func (r repository) Save(conn *gorm.DB, tStockmovement model.Stockmovement) error {
	return conn.Save(&tStockmovement).Error
}

func (r repository) Delete(conn *gorm.DB, tStockmovement model.Stockmovement) error {
	return conn.Delete(&tStockmovement).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageStockmovement) (vStockmovements []model.StockmovementView, count int64, err error) {
	query := conn.Model(&vStockmovements)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	for _, preload := range strings.Split(req.Preloads, ",") {
		conn = conn.Preload(preload)
	}

	if req.FromWarehouseID != "" {
		query = query.Where("from_warehouse_id = ?", req.FromWarehouseID)
	}
	if req.ToWarehouseID != "" {
		query = query.Where("to_warehouse_id = ?", req.ToWarehouseID)
	}
	if req.ProductID != "" {
		query = query.Where("product_id = ?", req.ProductID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Remark != "" {
		query = query.Where("remark ILIKE ?", "%"+req.Remark+"%")
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vStockmovements, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vStockmovements).Error
	if err != nil {
		return vStockmovements, count, err
	}

	return vStockmovements, count, err
}

func NewRepository() Repository {
	return repository{}
}
