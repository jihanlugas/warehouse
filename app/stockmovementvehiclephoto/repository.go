package stockmovementvehiclephoto

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tStockmovementvehiclephoto model.Stockmovementvehiclephoto, err error)
	GetTableByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (tStockmovementvehiclephoto model.Stockmovementvehiclephoto, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockmovementvehiclephoto model.StockmovementvehiclephotoView, err error)
	GetViewByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (vStockmovementvehiclephoto model.StockmovementvehiclephotoView, err error)
	Create(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error
	Update(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error
	Save(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error
	Delete(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error
	Page(conn *gorm.DB, req request.PageStockmovementvehiclephoto) (vStockmovementvehiclephotos []model.StockmovementvehiclephotoView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "stockmovementvehiclephoto"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tStockmovementvehiclephoto model.Stockmovementvehiclephoto, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tStockmovementvehiclephoto).Error
	return tStockmovementvehiclephoto, err
}

func (r repository) GetTableByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (tStockmovementvehiclephoto model.Stockmovementvehiclephoto, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("warehouse_id = ? ", warehouseID).Where("product_id = ? ", productID).First(&tStockmovementvehiclephoto).Error
	return tStockmovementvehiclephoto, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockmovementvehiclephoto model.StockmovementvehiclephotoView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vStockmovementvehiclephoto).Error
	return vStockmovementvehiclephoto, err
}

func (r repository) GetViewByWarehouseIdAndProductId(conn *gorm.DB, warehouseID, productID string, preloads ...string) (vStockmovementvehiclephoto model.StockmovementvehiclephotoView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("warehouse_id = ? ", warehouseID).Where("product_id = ? ", productID).First(&vStockmovementvehiclephoto).Error
	return vStockmovementvehiclephoto, err
}

func (r repository) Create(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error {
	return conn.Create(&tStockmovementvehiclephoto).Error
}

func (r repository) Update(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error {
	return conn.Model(&tStockmovementvehiclephoto).Updates(&tStockmovementvehiclephoto).Error
}

func (r repository) Save(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error {
	return conn.Save(&tStockmovementvehiclephoto).Error
}

func (r repository) Delete(conn *gorm.DB, tStockmovementvehiclephoto model.Stockmovementvehiclephoto) error {
	return conn.Delete(&tStockmovementvehiclephoto).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageStockmovementvehiclephoto) (vStockmovementvehiclephotos []model.StockmovementvehiclephotoView, count int64, err error) {
	query := conn.Model(&vStockmovementvehiclephotos)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.StockmovementvehicleID != "" {
		query = query.Where("stockmovementvehicle_id = ?", req.StockmovementvehicleID)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vStockmovementvehiclephotos, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vStockmovementvehiclephotos).Error
	if err != nil {
		return vStockmovementvehiclephotos, count, err
	}

	return vStockmovementvehiclephotos, count, err
}

func NewRepository() Repository {
	return repository{}
}
