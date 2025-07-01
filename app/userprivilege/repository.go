package userprivilege

import (
	"fmt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tUserprivilege model.Userprivilege, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vUserprivilege model.UserprivilegeView, err error)
	Create(conn *gorm.DB, tUserprivilege model.Userprivilege) error
	Update(conn *gorm.DB, tUserprivilege model.Userprivilege) error
	Save(conn *gorm.DB, tUserprivilege model.Userprivilege) error
	Delete(conn *gorm.DB, tUserprivilege model.Userprivilege) error
	Page(conn *gorm.DB, req request.PageUserprivilege) (vUserprivileges []model.UserprivilegeView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "userprivilege"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tUserprivilege model.Userprivilege, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tUserprivilege).Error
	return tUserprivilege, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tUserprivilege model.Userprivilege, err error) {
	err = conn.Where("name = ? ", name).First(&tUserprivilege).Error
	return tUserprivilege, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vUserprivilege model.UserprivilegeView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vUserprivilege).Error
	return vUserprivilege, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vUserprivilege model.UserprivilegeView, err error) {
	err = conn.Where("name = ? ", name).First(&vUserprivilege).Error
	return vUserprivilege, err
}

func (r repository) Create(conn *gorm.DB, tUserprivilege model.Userprivilege) error {
	return conn.Create(&tUserprivilege).Error
}

func (r repository) Update(conn *gorm.DB, tUserprivilege model.Userprivilege) error {
	return conn.Model(&tUserprivilege).Updates(&tUserprivilege).Error
}

func (r repository) Save(conn *gorm.DB, tUserprivilege model.Userprivilege) error {
	return conn.Save(&tUserprivilege).Error
}

func (r repository) Delete(conn *gorm.DB, tUserprivilege model.Userprivilege) error {
	return conn.Delete(&tUserprivilege).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageUserprivilege) (vUserprivileges []model.UserprivilegeView, count int64, err error) {
	query := conn.Model(&vUserprivileges)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.CreateName != "" {
		query = query.Where("create_name ILIKE ?", "%"+req.CreateName+"%")
	}

	err = query.Count(&count).Error
	if err != nil {
		return vUserprivileges, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vUserprivileges).Error
	if err != nil {
		return vUserprivileges, count, err
	}

	return vUserprivileges, count, err
}

func NewRepository() Repository {
	return repository{}
}
