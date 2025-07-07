package outbound

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vOutbound model.OutboundView, err error)
	Page(conn *gorm.DB, req request.PageOutbound) (vOutbounds []model.OutboundView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "outbound"
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vOutbound model.OutboundView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vOutbound).Error
	return vOutbound, err
}

func (r repository) Page(conn *gorm.DB, req request.PageOutbound) (vOutbounds []model.OutboundView, count int64, err error) {
	query := conn.Model(&vOutbounds)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.StockmovementID != "" {
		query = query.Where("stockmovement_id = ?", req.StockmovementID)
	}
	if req.ProductID != "" {
		query = query.Where("product_id = ?", req.ProductID)
	}
	if req.VehicleID != "" {
		query = query.Where("vehicle_id = ?", req.VehicleID)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Remark != "" {
		query = query.Where("remark ILIKE ?", "%"+req.Remark+"%")
	}
	if req.StartSentGrossQuantity != nil {
		query = query.Where("sent_gross_quantity >= ?", req.StartSentGrossQuantity)
	}
	if req.StartSentTareQuantity != nil {
		query = query.Where("sent_tare_quantity >= ?", req.StartSentTareQuantity)
	}
	if req.StartSentNetQuantity != nil {
		query = query.Where("sent_net_quantity >= ?", req.StartSentNetQuantity)
	}
	if req.StartSentTime != nil {
		query = query.Where("sent_time >= ?", req.StartSentTime)
	}
	if req.StartRecivedGrossQuantity != nil {
		query = query.Where("recived_gross_quantity >= ?", req.StartRecivedGrossQuantity)
	}
	if req.StartRecivedTareQuantity != nil {
		query = query.Where("recived_tare_quantity >= ?", req.StartRecivedTareQuantity)
	}
	if req.StartRecivedNetQuantity != nil {
		query = query.Where("recived_net_quantity >= ?", req.StartRecivedNetQuantity)
	}
	if req.StartRecivedTime != nil {
		query = query.Where("recived_time >= ?", req.StartRecivedTime)
	}
	if req.EndSentGrossQuantity != nil {
		query = query.Where("sent_gross_quantity <= ?", req.EndSentGrossQuantity)
	}
	if req.EndSentTareQuantity != nil {
		query = query.Where("sent_tare_quantity <= ?", req.EndSentTareQuantity)
	}
	if req.EndSentNetQuantity != nil {
		query = query.Where("sent_net_quantity <= ?", req.EndSentNetQuantity)
	}
	if req.EndSentTime != nil {
		query = query.Where("sent_time <= ?", req.EndSentTime)
	}
	if req.EndRecivedGrossQuantity != nil {
		query = query.Where("recived_gross_quantity <= ?", req.EndRecivedGrossQuantity)
	}
	if req.EndRecivedTareQuantity != nil {
		query = query.Where("recived_tare_quantity <= ?", req.EndRecivedTareQuantity)
	}
	if req.EndRecivedNetQuantity != nil {
		query = query.Where("recived_net_quantity <= ?", req.EndRecivedNetQuantity)
	}
	if req.EndRecivedTime != nil {
		query = query.Where("recived_time <= ?", req.EndRecivedTime)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vOutbounds, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vOutbounds).Error
	if err != nil {
		return vOutbounds, count, err
	}

	return vOutbounds, count, err
}

func NewRepository() Repository {
	return repository{}
}
