package stockin

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockin model.StockinView, err error)
	Page(conn *gorm.DB, req request.PageStockin) (vStockins []model.StockinView, count int64, err error)
}

type repository struct {
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockin model.StockinView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vStockin).Error
	return vStockin, err
}

func (r repository) Page(conn *gorm.DB, req request.PageStockin) (vStockins []model.StockinView, count int64, err error) {
	query := conn.Model(&vStockins)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID != "" {
		query = query.Where("product_id = ?", req.ProductID)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.Remark != "" {
		query = query.Where("remark ILIKE ?", "%"+req.Remark+"%")
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
		return vStockins, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vStockins).Error
	if err != nil {
		return vStockins, count, err
	}

	return vStockins, count, err
}

func NewRepository() Repository {
	return repository{}
}
