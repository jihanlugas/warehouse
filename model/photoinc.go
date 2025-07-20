package model

import (
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"time"
)

func (m *Photoinc) BeforeCreate(tx *gorm.DB) (err error) {
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

func (m *Photoinc) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *PhotoincView) AfterFind(tx *gorm.DB) (err error) {
	return
}
