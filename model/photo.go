package model

import (
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"time"
)

func (m *Photo) BeforeCreate(tx *gorm.DB) (err error) {
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

func (m *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateDt = now
	return
}

func (m *PhotoView) AfterFind(tx *gorm.DB) (err error) {
	return
}
