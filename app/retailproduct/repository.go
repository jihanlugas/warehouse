package retailproduct

import (
	"github.com/jihanlugas/warehouse/model"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tRetailproduct model.Retailproduct, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vRetailproduct model.RetailproductView, err error)
	Create(conn *gorm.DB, tRetailproduct model.Retailproduct) error
	Update(conn *gorm.DB, tRetailproduct model.Retailproduct) error
	Save(conn *gorm.DB, tRetailproduct model.Retailproduct) error
	Delete(conn *gorm.DB, tRetailproduct model.Retailproduct) error
}

type repository struct {
}

func (r repository) Name() string {
	return "retailproduct"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tRetailproduct model.Retailproduct, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tRetailproduct).Error
	return tRetailproduct, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tRetailproduct model.Retailproduct, err error) {
	err = conn.Where("name = ? ", name).First(&tRetailproduct).Error
	return tRetailproduct, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vRetailproduct model.RetailproductView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vRetailproduct).Error
	return vRetailproduct, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vRetailproduct model.RetailproductView, err error) {
	err = conn.Where("name = ? ", name).First(&vRetailproduct).Error
	return vRetailproduct, err
}

func (r repository) Create(conn *gorm.DB, tRetailproduct model.Retailproduct) error {
	return conn.Create(&tRetailproduct).Error
}

func (r repository) Update(conn *gorm.DB, tRetailproduct model.Retailproduct) error {
	return conn.Model(&tRetailproduct).Updates(&tRetailproduct).Error
}

func (r repository) Save(conn *gorm.DB, tRetailproduct model.Retailproduct) error {
	return conn.Save(&tRetailproduct).Error
}

func (r repository) Delete(conn *gorm.DB, tRetailproduct model.Retailproduct) error {
	return conn.Delete(&tRetailproduct).Error
}

func NewRepository() Repository {
	return repository{}
}
