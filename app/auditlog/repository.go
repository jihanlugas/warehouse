package auditlog

import (
	"fmt"
	"strings"

	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"gorm.io/gorm"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tAuditlog model.Auditlog, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vAuditlog model.AuditlogView, err error)
	Create(conn *gorm.DB, tAuditlog model.Auditlog) error
	Update(conn *gorm.DB, tAuditlog model.Auditlog) error
	Save(conn *gorm.DB, tAuditlog model.Auditlog) error
	Delete(conn *gorm.DB, tAuditlog model.Auditlog) error
	Page(conn *gorm.DB, req request.PageAuditlog) (vAuditlogs []model.AuditlogView, count int64, err error)
}

type repository struct {
}

func (r repository) Name() string {
	return "auditlog"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tAuditlog model.Auditlog, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tAuditlog).Error
	return tAuditlog, err
}

func (r repository) GetTableByName(conn *gorm.DB, name string) (tAuditlog model.Auditlog, err error) {
	err = conn.Where("name = ? ", name).First(&tAuditlog).Error
	return tAuditlog, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vAuditlog model.AuditlogView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vAuditlog).Error
	return vAuditlog, err
}

func (r repository) GetViewByName(conn *gorm.DB, name string) (vAuditlog model.AuditlogView, err error) {
	err = conn.Where("name = ? ", name).First(&vAuditlog).Error
	return vAuditlog, err
}

func (r repository) Create(conn *gorm.DB, tAuditlog model.Auditlog) error {
	return conn.Create(&tAuditlog).Error
}

func (r repository) Update(conn *gorm.DB, tAuditlog model.Auditlog) error {
	return conn.Model(&tAuditlog).Updates(&tAuditlog).Error
}

func (r repository) Save(conn *gorm.DB, tAuditlog model.Auditlog) error {
	return conn.Save(&tAuditlog).Error
}

func (r repository) Delete(conn *gorm.DB, tAuditlog model.Auditlog) error {
	return conn.Delete(&tAuditlog).Error
}

func (r repository) Page(conn *gorm.DB, req request.PageAuditlog) (vAuditlogs []model.AuditlogView, count int64, err error) {
	query := conn.Model(&vAuditlogs)

	if req.Preloads != "" {
		preloads := strings.Split(req.Preloads, ",")
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	if req.LocationID != "" {
		query = query.Where("location_id = ?", req.LocationID)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.AuditlogType != "" {
		query = query.Where("auditlog_type = ?", req.AuditlogType)
	}
	if req.Title != "" {
		query = query.Where("title ILIKE ?", "%"+req.Title+"%")
	}
	if req.Description != "" {
		query = query.Where("description ILIKE ?", "%"+req.Description+"%")
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
		return vAuditlogs, count, err
	}

	if req.SortField != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, req.SortOrder))
	} else {
		query = query.Order(fmt.Sprintf("%s %s", "create_dt", "desc"))
	}
	if req.Limit >= 0 {
		query = query.Offset((req.GetPage() - 1) * req.GetLimit()).Limit(req.GetLimit())
	}

	err = query.Find(&vAuditlogs).Error
	if err != nil {
		return vAuditlogs, count, err
	}

	return vAuditlogs, count, err
}

func NewRepository() Repository {
	return repository{}
}
