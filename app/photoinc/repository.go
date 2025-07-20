package photoinc

import (
	"errors"
	"github.com/jihanlugas/warehouse/config"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"strconv"
)

type Repository interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPhotoinc model.Photoinc, err error)
	GetTableToUse(conn *gorm.DB, refTable model.PhotoRef) (tPhotoinc model.Photoinc, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPhotoinc model.PhotoincView, err error)
	AddRunning(conn *gorm.DB, tPhotoinc model.Photoinc) error
	Create(conn *gorm.DB, tPhotoinc model.Photoinc) error
	Update(conn *gorm.DB, tPhotoinc model.Photoinc) error
	Save(conn *gorm.DB, tPhotoinc model.Photoinc) error
	Delete(conn *gorm.DB, tPhotoinc model.Photoinc) error
}

type repository struct {
}

func (r repository) Name() string {
	return "photoinc"
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPhotoinc model.Photoinc, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tPhotoinc).Error
	return tPhotoinc, err
}

func (r repository) GetTableToUse(conn *gorm.DB, refTable model.PhotoRef) (tPhotoinc model.Photoinc, err error) {
	err = conn.Where("ref_table = ? ", refTable).
		Order("ref_table DESC").
		First(&tPhotoinc).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			tPhotoinc.RefTable = string(refTable)
			tPhotoinc.FolderInc = 1
			tPhotoinc.Folder = config.StorageDirectory + "/" + config.PhotoDirectory + "/" + tPhotoinc.RefTable + "/" + strconv.FormatInt(tPhotoinc.FolderInc, 10)
			tPhotoinc.Running = 0
			err = r.Create(conn, tPhotoinc)
			if err != nil {
				return tPhotoinc, err
			}
			err = utils.CreateFolder(tPhotoinc.Folder, 0755)
			if err != nil {
				return tPhotoinc, err
			}
		} else {
			return tPhotoinc, err
		}
	} else {
		if tPhotoinc.Running >= config.PhotoincRunningLimit {
			tPhotoinc.RefTable = (string(refTable))
			tPhotoinc.FolderInc = tPhotoinc.FolderInc + 1
			tPhotoinc.Folder = config.StorageDirectory + "/" + config.PhotoDirectory + "/" + tPhotoinc.RefTable + "/" + strconv.FormatInt(tPhotoinc.FolderInc, 10)
			tPhotoinc.Running = 0
			err = r.Create(conn, tPhotoinc)
			if err != nil {
				return tPhotoinc, err
			}

			err = utils.CreateFolder(tPhotoinc.Folder, 0755)
			if err != nil {
				return tPhotoinc, err
			}
		}
	}
	return tPhotoinc, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPhotoinc model.PhotoincView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vPhotoinc).Error
	return vPhotoinc, err
}

func (r repository) AddRunning(conn *gorm.DB, tPhotoinc model.Photoinc) (err error) {
	tPhotoinc.Running++
	return conn.Save(&tPhotoinc).Error
}

func (r repository) Create(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	return conn.Create(&tPhotoinc).Error
}

func (r repository) Update(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	return conn.Model(&tPhotoinc).Updates(&tPhotoinc).Error
}

func (r repository) Save(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	return conn.Save(&tPhotoinc).Error
}

func (r repository) Delete(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	return conn.Delete(&tPhotoinc).Error
}

func NewRepository() Repository {
	return repository{}
}
