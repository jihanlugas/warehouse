package customer

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tCustomer model.Customer, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vCustomer model.CustomerView, err error)
	Create(conn *gorm.DB, tCustomer model.Customer) error
	Update(conn *gorm.DB, tCustomer model.Customer) error
	Save(conn *gorm.DB, tCustomer model.Customer) error
	Delete(conn *gorm.DB, tCustomer model.Customer) error
	Page(conn *gorm.DB, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "customer"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tCustomer model.Customer, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tCustomer).Error
	return tCustomer, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tCustomer model.Customer, err error) {
	err = conn.Where("name = ? ", name).First(&tCustomer).Error
	return tCustomer, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vCustomer model.CustomerView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vCustomer).Error
	return vCustomer, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vCustomer model.CustomerView, err error) {
	err = conn.Where("name = ? ", name).First(&vCustomer).Error
	return vCustomer, err
}

func (r repository) Create(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Create(&tCustomer).Error
}

func (r repository) Update(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Model(&tCustomer).Updates(&tCustomer).Error
}

func (r repository) Save(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Save(&tCustomer).Error
}

func (r repository) Delete(conn *gorm.DB, tCustomer model.Customer) error {
	return conn.Delete(&tCustomer).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageCustomer) (vCustomers []model.CustomerView, count int64, err error) {
	query := conn.Model(&vCustomers)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.PhoneNumber != "" {
		query = query.Where("phoneNumber ILIKE ?", "%"+utils.FormatPhoneTo62(req.PhoneNumber)+"%")
	}
	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}
	if req.Address != "" {
		query = query.Where("address ILIKE ?", "%"+req.Address+"%")
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
		return vCustomers, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vCustomers).Error
	if err != nil {
		return vCustomers, count, err
	}

	return vCustomers, count, err
}

func NewRepository() Repository {
	return repository{}
}
