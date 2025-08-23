package purchaseorder

import (
	"fmt"
	"strings"
	"time"

	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPurchaseorder model.Purchaseorder, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPurchaseorder model.PurchaseorderView, err error)
	GetNextNumber(conn *gorm.DB) (number int64)
	Create(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error
	Update(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error
	Save(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error
	Delete(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error
	Page(conn *gorm.DB, req request.PagePurchaseorder) (vPurchaseorders []model.PurchaseorderView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "purchaseorder"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPurchaseorder model.Purchaseorder, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tPurchaseorder).Error
	return tPurchaseorder, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tPurchaseorder model.Purchaseorder, err error) {
	err = conn.Where("name = ? ", name).First(&tPurchaseorder).Error
	return tPurchaseorder, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPurchaseorder model.PurchaseorderView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vPurchaseorder).Error
	return vPurchaseorder, err
}

func (r repository) GetNextNumber(conn *gorm.DB) (number int64) {
	conn.Model(&model.Purchaseorder{}).Unscoped().
		Where("EXTRACT(MONTH FROM create_dt) = EXTRACT(MONTH FROM CURRENT_DATE)").
		Where("EXTRACT(YEAR FROM create_dt) = EXTRACT(YEAR FROM CURRENT_DATE)").
		Count(&number)
	return number + 1
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vPurchaseorder model.PurchaseorderView, err error) {
	err = conn.Where("name = ? ", name).First(&vPurchaseorder).Error
	return vPurchaseorder, err
}

func (r repository) Create(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error {
	var date time.Time

	if tPurchaseorder.CreateDt.IsZero() {
		date = time.Now()
	} else {
		date = tPurchaseorder.CreateDt
	}
	tPurchaseorder.Number = fmt.Sprintf("%d/PO/%s/%d", r.GetNextNumber(conn), utils.DisplayRoman(int(date.Month())), date.Year())
	return conn.Create(&tPurchaseorder).Error
}

func (r repository) Update(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error {
	return conn.Model(&tPurchaseorder).Updates(&tPurchaseorder).Error
}

func (r repository) Save(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error {
	return conn.Save(&tPurchaseorder).Error
}

func (r repository) Delete(conn *gorm.DB, tPurchaseorder model.Purchaseorder) error {
	return conn.Delete(&tPurchaseorder).Error
}

func (r repository) Page(conn *gorm.DB, req request.PagePurchaseorder) (vPurchaseorders []model.PurchaseorderView, count int64, err error) {
	query := conn.Model(&vPurchaseorders)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.Notes != "" {
		query = query.Where("notes ILIKE ?", "%"+req.Notes+"%")
	}
	if req.Number != "" {
		query = query.Where("number ILIKE ?", "%"+req.Number+"%")
	}
	if req.PurchaseorderStatus != nil {
		query = query.Where("purchaseorder_status = ?", req.PurchaseorderStatus)
	}
	if req.StartTotalPrice != nil {
		query = query.Where("total_price >= ?", req.StartTotalPrice)
	}
	if req.EndTotalPrice != nil {
		query = query.Where("total_price <= ?", req.EndTotalPrice)
	}
	if req.StartTotalPayment != nil {
		query = query.Where("total_payment >= ?", req.StartTotalPayment)
	}
	if req.EndTotalPayment != nil {
		query = query.Where("total_payment <= ?", req.EndTotalPayment)
	}
	if req.StartOutstanding != nil {
		query = query.Where("outstanding >= ?", req.StartOutstanding)
	}
	if req.EndOutstanding != nil {
		query = query.Where("outstanding <= ?", req.EndOutstanding)
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
		return vPurchaseorders, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vPurchaseorders).Error
	if err != nil {
		return vPurchaseorders, count, err
	}

	return vPurchaseorders, count, err
}

func NewRepository() Repository {
	return repository{}
}
