package stockmovementvehicle

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tStockmovementvehicle model.Stockmovementvehicle, err error)
	GetTableByName(conn *gorm.DB, name string) (tStockmovementvehicle model.Stockmovementvehicle, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	GetNextNumber(conn *gorm.DB) (number int64)
	GetViewByName(conn *gorm.DB, name string) (vStockmovementvehicle model.StockmovementvehicleView, err error)
	Create(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error
	Update(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error
	Save(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error
	Delete(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error
	Page(conn *gorm.DB, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "stockmovementvehicle"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tStockmovementvehicle model.Stockmovementvehicle, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tStockmovementvehicle).Error
	return tStockmovementvehicle, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tStockmovementvehicle model.Stockmovementvehicle, err error) {
	err = conn.Where("name = ? ", name).First(&tStockmovementvehicle).Error
	return tStockmovementvehicle, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vStockmovementvehicle model.StockmovementvehicleView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vStockmovementvehicle).Error
	return vStockmovementvehicle, err
}

func (r repository) GetNextNumber(conn *gorm.DB) (number int64) {
	conn.Model(&model.Stockmovementvehicle{}).Unscoped().
		Where("EXTRACT(MONTH FROM create_dt) = EXTRACT(MONTH FROM CURRENT_DATE)").
		Where("EXTRACT(YEAR FROM create_dt) = EXTRACT(YEAR FROM CURRENT_DATE)").
		Count(&number)
	return number + 1
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vStockmovementvehicle model.StockmovementvehicleView, err error) {
	err = conn.Where("name = ? ", name).First(&vStockmovementvehicle).Error
	return vStockmovementvehicle, err
}

func (r repository) Create(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error {
	var date time.Time

	if tStockmovementvehicle.CreateDt.IsZero() {
		date = time.Now()
	} else {
		date = tStockmovementvehicle.CreateDt
	}
	tStockmovementvehicle.Number = fmt.Sprintf("%d/DELIVERY/%s/%d", r.GetNextNumber(conn), utils.DisplayRoman(int(date.Month())), date.Year())
	return conn.Create(&tStockmovementvehicle).Error
}

func (r repository) Update(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error {
	return conn.Model(&tStockmovementvehicle).Updates(&tStockmovementvehicle).Error
}

func (r repository) Save(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error {
	return conn.Save(&tStockmovementvehicle).Error
}

func (r repository) Delete(conn *gorm.DB, tStockmovementvehicle model.Stockmovementvehicle) error {
	return conn.Delete(&tStockmovementvehicle).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageStockmovementvehicle) (vStockmovementvehicles []model.StockmovementvehicleView, count int64, err error) {
	query := conn.Model(&vStockmovementvehicles)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.FromWarehouseID != "" {
		query = query.Where("from_warehouse_id = ?", req.FromWarehouseID)
	}
	if req.ToWarehouseID != "" {
		query = query.Where("to_warehouse_id = ?", req.ToWarehouseID)
	}
	if req.StockmovementID != "" {
		query = query.Where("stockmovement_id = ?", req.StockmovementID)
	}
	if req.ProductID != "" {
		query = query.Where("product_id = ?", req.ProductID)
	}
	if req.RelatedID != "" {
		query = query.Where("related_id = ?", req.RelatedID)
	}
	if req.VehicleID != "" {
		query = query.Where("vehicle_id = ?", req.VehicleID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
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
		return vStockmovementvehicles, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vStockmovementvehicles).Error
	if err != nil {
		return vStockmovementvehicles, count, err
	}

	return vStockmovementvehicles, count, err
}

func NewRepository() Repository {
	return repository{}
}
