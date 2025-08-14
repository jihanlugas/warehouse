package transaction

import (
	"fmt"
	"strings"

	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tTransaction model.Transaction, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vTransaction model.TransactionView, err error)
	Create(conn *gorm.DB, tTransaction model.Transaction) error
	Update(conn *gorm.DB, tTransaction model.Transaction) error
	Save(conn *gorm.DB, tTransaction model.Transaction) error
	Delete(conn *gorm.DB, tTransaction model.Transaction) error
	Page(conn *gorm.DB, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "transaction"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tTransaction model.Transaction, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tTransaction).Error
	return tTransaction, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tTransaction model.Transaction, err error) {
	err = conn.Where("name = ? ", name).First(&tTransaction).Error
	return tTransaction, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vTransaction model.TransactionView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vTransaction).Error
	return vTransaction, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vTransaction model.TransactionView, err error) {
	err = conn.Where("name = ? ", name).First(&vTransaction).Error
	return vTransaction, err
}

func (r repository) Create(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Create(&tTransaction).Error
}

func (r repository) Update(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Model(&tTransaction).Updates(&tTransaction).Error
}

func (r repository) Save(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Save(&tTransaction).Error
}

func (r repository) Delete(conn *gorm.DB, tTransaction model.Transaction) error {
	return conn.Delete(&tTransaction).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageTransaction) (vTransactions []model.TransactionView, count int64, err error) {
	query := conn.Model(&vTransactions)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.RelatedID != "" {
		query = query.Where("related_id = ?", req.RelatedID)
	}
	if req.TransactionRelated != "" {
		query = query.Where("transaction_related = ?", req.TransactionRelated)
	}
	if req.Notes != "" {
		query = query.Where("notes ILIKE ?", "%"+req.Notes+"%")
	}
	if req.StartAmount != nil {
		query = query.Where("amount >= ?", req.StartAmount)
	}
	if req.EndAmount != nil {
		query = query.Where("amount <= ?", req.EndAmount)
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
		return vTransactions, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vTransactions).Error
	if err != nil {
		return vTransactions, count, err
	}

	return vTransactions, count, err
}

func NewRepository() Repository {
	return repository{}
}
