package stocklog

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tStocklog model.Stocklog, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vStocklog model.StocklogView, err error)
	Create(conn *gorm.DB, tStocklog model.Stocklog) error
	Update(conn *gorm.DB, tStocklog model.Stocklog) error
	Save(conn *gorm.DB, tStocklog model.Stocklog) error
	Delete(conn *gorm.DB, tStocklog model.Stocklog) error
	Page(conn *gorm.DB, req request.PageStocklog) (vStocklogs []model.StocklogView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "stocklog"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tStocklog model.Stocklog, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tStocklog).Error
	return tStocklog, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vStocklog model.StocklogView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vStocklog).Error
	return vStocklog, err
}

func (r repository) Create(conn *gorm.DB, tStocklog model.Stocklog) error {
	return conn.Create(&tStocklog).Error
}

func (r repository) Update(conn *gorm.DB, tStocklog model.Stocklog) error {
	return conn.Model(&tStocklog).Updates(&tStocklog).Error
}

func (r repository) Save(conn *gorm.DB, tStocklog model.Stocklog) error {
	return conn.Save(&tStocklog).Error
}

func (r repository) Delete(conn *gorm.DB, tStocklog model.Stocklog) error {
	return conn.Delete(&tStocklog).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageStocklog) (vStocklogs []model.StocklogView, count int64, err error) {
	query := conn.Model(&vStocklogs)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.StockID != "" {
		query = query.Where("stock_id = ?", req.StockID)
	}
	if req.StockmovementID != "" {
		query = query.Where("stockmovement_id = ?", req.StockmovementID)
	}
	if req.StockmovementvehicleID != "" {
		query = query.Where("stockmovementvehicle_id = ?", req.StockmovementvehicleID)
	}
	if req.ProductID != "" {
		query = query.Where("product_id = ?", req.ProductID)
	}
	if req.VehicleID != "" {
		query = query.Where("vehicle_id = ?", req.VehicleID)
	}
	if req.StartGrossQuantity != nil {
		query = query.Where("gross_quantity >= ?", req.StartGrossQuantity)
	}
	if req.StartTareQuantity != nil {
		query = query.Where("tare_quantity >= ?", req.StartTareQuantity)
	}
	if req.StartNetQuantity != nil {
		query = query.Where("net_quantity >= ?", req.StartNetQuantity)
	}
	if req.StartCreateDt != nil {
		query = query.Where("create_dt >= ?", req.StartCreateDt)
	}
	if req.EndGrossQuantity != nil {
		query = query.Where("gross_quantity <= ?", req.EndGrossQuantity)
	}
	if req.EndTareQuantity != nil {
		query = query.Where("tare_quantity <= ?", req.EndTareQuantity)
	}
	if req.EndNetQuantity != nil {
		query = query.Where("net_quantity <= ?", req.EndNetQuantity)
	}
	if req.EndCreateDt != nil {
		query = query.Where("create_dt <= ?", req.EndCreateDt)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vStocklogs, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vStocklogs).Error
	if err != nil {
		return vStocklogs, count, err
	}

	return vStocklogs, count, err
}

func NewRepository() Repository {
	return repository{}
}
