package retail

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
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tRetail model.Retail, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vRetail model.RetailView, err error)
	GetNextNumber(conn *gorm.DB) (number int64)
	Create(conn *gorm.DB, tRetail model.Retail) error
	Update(conn *gorm.DB, tRetail model.Retail) error
	Save(conn *gorm.DB, tRetail model.Retail) error
	Delete(conn *gorm.DB, tRetail model.Retail) error
	Page(conn *gorm.DB, req request.PageRetail) (vRetails []model.RetailView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "retail"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tRetail model.Retail, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tRetail).Error
	return tRetail, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tRetail model.Retail, err error) {
	err = conn.Where("name = ? ", name).First(&tRetail).Error
	return tRetail, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vRetail model.RetailView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vRetail).Error
	return vRetail, err
}

func (r repository) GetNextNumber(conn *gorm.DB) (number int64) {
	conn.Model(&model.Retail{}).Unscoped().
		Where("EXTRACT(MONTH FROM create_dt) = EXTRACT(MONTH FROM CURRENT_DATE)").
		Where("EXTRACT(YEAR FROM create_dt) = EXTRACT(YEAR FROM CURRENT_DATE)").
		Count(&number)
	return number + 1
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vRetail model.RetailView, err error) {
	err = conn.Where("name = ? ", name).First(&vRetail).Error
	return vRetail, err
}

func (r repository) Create(conn *gorm.DB, tRetail model.Retail) error {
	var date time.Time

	if tRetail.CreateDt.IsZero() {
		date = time.Now()
	} else {
		date = tRetail.CreateDt
	}
	tRetail.Number = fmt.Sprintf("%d/RETAIL/%s/%d", r.GetNextNumber(conn), utils.DisplayRoman(int(date.Month())), date.Year())
	return conn.Create(&tRetail).Error
}

func (r repository) Update(conn *gorm.DB, tRetail model.Retail) error {
	return conn.Model(&tRetail).Updates(&tRetail).Error
}

func (r repository) Save(conn *gorm.DB, tRetail model.Retail) error {
	return conn.Save(&tRetail).Error
}

func (r repository) Delete(conn *gorm.DB, tRetail model.Retail) error {
	return conn.Delete(&tRetail).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageRetail) (vRetails []model.RetailView, count int64, err error) {
	query := conn.Model(&vRetails)

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
	if req.Status != nil {
		query = query.Where("status = ?", req.Status)
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
		return vRetails, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vRetails).Error
	if err != nil {
		return vRetails, count, err
	}

	return vRetails, count, err
}

func NewRepository() Repository {
	return repository{}
}
