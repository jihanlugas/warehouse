package model

import (
	"time"

	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
)

func (m *Auditlog) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()

	if m.ID == "" {
		m.ID = utils.GetUniqueID()
	}

	if m.CreateDt.IsZero() {
		m.CreateDt = now
	}
	if m.UpdateDt.IsZero() {
		m.UpdateDt = now
	}
	return
}

func (m *Auditlog) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *AuditlogView) AfterFind(tx *gorm.DB) (err error) {
	return
}
