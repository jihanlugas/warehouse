package purchaseorderproduct

import (
	"github.com/jihanlugas/warehouse/model"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPurchaseorderproduct model.Purchaseorderproduct, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPurchaseorderproduct model.PurchaseorderproductView, err error)
	Create(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error
	Update(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error
	Save(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error
	Delete(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error
}

type repository struct {
}

func (r repository) Name() string {
	return "purchaseorderproduct"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPurchaseorderproduct model.Purchaseorderproduct, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tPurchaseorderproduct).Error
	return tPurchaseorderproduct, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tPurchaseorderproduct model.Purchaseorderproduct, err error) {
	err = conn.Where("name = ? ", name).First(&tPurchaseorderproduct).Error
	return tPurchaseorderproduct, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPurchaseorderproduct model.PurchaseorderproductView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vPurchaseorderproduct).Error
	return vPurchaseorderproduct, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vPurchaseorderproduct model.PurchaseorderproductView, err error) {
	err = conn.Where("name = ? ", name).First(&vPurchaseorderproduct).Error
	return vPurchaseorderproduct, err
}

func (r repository) Create(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error {
	return conn.Create(&tPurchaseorderproduct).Error
}

func (r repository) Update(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error {
	return conn.Model(&tPurchaseorderproduct).Updates(&tPurchaseorderproduct).Error
}

func (r repository) Save(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error {
	return conn.Save(&tPurchaseorderproduct).Error
}

func (r repository) Delete(conn *gorm.DB, tPurchaseorderproduct model.Purchaseorderproduct) error {
	return conn.Delete(&tPurchaseorderproduct).Error
}

func NewRepository() Repository {
	return repository{}
}
