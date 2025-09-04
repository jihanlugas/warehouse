package userprovider

import (
	"github.com/jihanlugas/warehouse/model"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tUserprovider model.Userprovider, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vUserprovider model.UserproviderView, err error)
	GetTableByProviderNameAndUserId(conn *gorm.DB, providerName, userId string, preloads ...string) (tUserprovider model.Userprovider, err error)
	GetViewByProviderNameAndUserId(conn *gorm.DB, providerName, userId string, preloads ...string) (vUserprovider model.UserproviderView, err error)
	GetTableByProviderNameAndEmail(conn *gorm.DB, providerName string, email string, preloads ...string) (tUserprovider model.Userprovider, err error)
	GetViewByProviderNameAndEmail(conn *gorm.DB, providerName string, email string, preloads ...string) (vUserprovider model.UserproviderView, err error)
	Create(conn *gorm.DB, tUserprovider model.Userprovider) error
	Update(conn *gorm.DB, tUserprovider model.Userprovider) error
	Save(conn *gorm.DB, tUserprovider model.Userprovider) error
	Delete(conn *gorm.DB, tUserprovider model.Userprovider) error
}

type repository struct {
}

func (r repository) Name() string {
	return "userprovider"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tUserprovider model.Userprovider, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tUserprovider).Error
	return tUserprovider, err
}

func (r repository) GetViewById(conn *gorm.DB, providerUserId string, preloads ...string) (vUserprovider model.UserproviderView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("provider_user_id = ? ", providerUserId).First(&vUserprovider).Error
	return vUserprovider, err
}

func (r repository) GetTableByProviderNameAndUserId(conn *gorm.DB, providerName, userId string, preloads ...string) (tUserprovider model.Userprovider, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("provider_name = ? ", providerName).
		Where("user_id = ? ", userId).
		First(&tUserprovider).Error
	return tUserprovider, err
}

func (r repository) GetViewByProviderNameAndUserId(conn *gorm.DB, providerName, userId string, preloads ...string) (vUserprovider model.UserproviderView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("provider_name = ? ", providerName).
		Where("user_id = ? ", userId).
		First(&vUserprovider).Error
	return vUserprovider, err
}

func (r repository) GetTableByProviderNameAndEmail(conn *gorm.DB, providerName string, email string, preloads ...string) (tUserprovider model.Userprovider, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("provider_name = ? ", providerName).
		Where("email = ? ", email).
		First(&tUserprovider).Error
	return tUserprovider, err
}

func (r repository) GetViewByProviderNameAndEmail(conn *gorm.DB, providerName string, email string, preloads ...string) (vUserprovider model.UserproviderView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("provider_name = ? ", providerName).
		Where("email = ? ", email).
		First(&vUserprovider).Error
	return vUserprovider, err
}

func (r repository) Create(conn *gorm.DB, tUserprovider model.Userprovider) error {
	return conn.Create(&tUserprovider).Error
}

func (r repository) Update(conn *gorm.DB, tUserprovider model.Userprovider) error {
	return conn.Model(&tUserprovider).Updates(&tUserprovider).Error
}

func (r repository) Save(conn *gorm.DB, tUserprovider model.Userprovider) error {
	return conn.Save(&tUserprovider).Error
}

func (r repository) Delete(conn *gorm.DB, tUserprovider model.Userprovider) error {
	return conn.Delete(&tUserprovider).Error
}

func NewRepository() Repository {
	return repository{}
}
