package user

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
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tUser model.User, err error)
	GetByUsername(conn *gorm.DB, username string, preloads ...string) (tUser model.User, err error)
	GetByEmail(conn *gorm.DB, email string, preloads ...string) (tUser model.User, err error)
	GetByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (tUser model.User, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vUser model.UserView, err error)
	GetViewByUsername(conn *gorm.DB, username string, preloads ...string) (vUser model.UserView, err error)
	GetViewByEmail(conn *gorm.DB, email string, preloads ...string) (vUser model.UserView, err error)
	GetViewByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (vUser model.UserView, err error)
	Create(conn *gorm.DB, tUser model.User) error
	Update(conn *gorm.DB, tUser model.User) error
	Save(conn *gorm.DB, tUser model.User) error
	Delete(conn *gorm.DB, tUser model.User) error
	Page(conn *gorm.DB, req request.PageUser) (vUsers []model.UserView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "user"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByUsername(conn *gorm.DB, username string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("username = ? ", username).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByEmail(conn *gorm.DB, email string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("email = ? ", email).First(&tUser).Error
	return tUser, err
}

func (r repository) GetByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (tUser model.User, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("phone_number = ? ", utils.FormatPhoneTo62(phoneNumber)).First(&tUser).Error
	return tUser, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("id = ? ", id).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByUsername(conn *gorm.DB, username string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("username = ? ", username).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByEmail(conn *gorm.DB, email string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("email = ? ", email).First(&vUser).Error
	return vUser, err
}

func (r repository) GetViewByPhoneNumber(conn *gorm.DB, phoneNumber string, preloads ...string) (vUser model.UserView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	err = conn.Where("phone_number = ? ", phoneNumber).First(&vUser).Error
	return vUser, err
}

func (r repository) Create(conn *gorm.DB, tUser model.User) error {
	return conn.Create(&tUser).Error
}

func (r repository) Update(conn *gorm.DB, tUser model.User) error {
	return conn.Model(&tUser).Updates(&tUser).Error
}

func (r repository) Save(conn *gorm.DB, tUser model.User) error {
	return conn.Save(&tUser).Error
}

func (r repository) Delete(conn *gorm.DB, tUser model.User) error {
	return conn.Delete(&tUser).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageUser) (vUsers []model.UserView, count int64, err error) {
	query := conn.Model(&vUsers)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	// query
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}
	if req.Fullname != "" {
		query = query.Where("fullname ILIKE ?", "%"+req.Fullname+"%")
	}
	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+req.Username+"%")
	}
	if req.PhoneNumber != "" {
		query = query.Where("phone_number ILIKE ?", "%"+utils.FormatPhoneTo62(req.PhoneNumber)+"%")
	}
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+utils.FormatPhoneTo62(req.Username)+"%")
	}
	if req.Address != "" {
		query = query.Where("address ILIKE ?", "%"+utils.FormatPhoneTo62(req.Address)+"%")
	}
	if req.BirthPlace != "" {
		query = query.Where("birth_place ILIKE ?", "%"+utils.FormatPhoneTo62(req.BirthPlace)+"%")
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
		return vUsers, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "fullname", "asc"))
	}

	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vUsers).Error
	if err != nil {
		return vUsers, count, err
	}

	return vUsers, count, err
}

func NewRepository() Repository {
	return repository{}
}
