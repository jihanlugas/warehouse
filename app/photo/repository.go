package photo

import (
	"errors"
	"fmt"
	"github.com/jihanlugas/warehouse/config"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
)

type Repository interface {
	Name() string
	Upload(conn *gorm.DB, file *multipart.FileHeader, photoRef model.PhotoRef) (tPhoto model.Photo, err error)
	GetTableById(conn *gorm.DB, id string, preloads ...string) (tPhoto model.Photo, err error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (vPhoto model.PhotoView, err error)
	Create(conn *gorm.DB, tPhoto model.Photo) error
	Update(conn *gorm.DB, tPhoto model.Photo) error
	Save(conn *gorm.DB, tPhoto model.Photo) error
	Delete(conn *gorm.DB, tPhoto model.Photo) error
}

type repository struct {
}

func (r repository) Name() string {
	return "photo"
}

func (r repository) Upload(conn *gorm.DB, file *multipart.FileHeader, photoRef model.PhotoRef) (tPhoto model.Photo, err error) {
	tPhotoinc, err := r.photoincGettouse(conn, photoRef)
	if err != nil {
		return tPhoto, err
	}

	tPhoto.ID = utils.GetUniqueID()
	tPhoto.Ext = filepath.Ext(file.Filename)
	tPhoto.ClientName = file.Filename
	tPhoto.ServerName = fmt.Sprintf("%s%s", tPhoto.ID, tPhoto.Ext)
	tPhoto.RefTable = string(photoRef)
	tPhoto.PhotoPath = fmt.Sprintf("%s/%s", tPhotoinc.Folder, tPhoto.ServerName)
	tPhoto.PhotoSize = file.Size
	tPhoto.PhotoWidth = 0
	tPhoto.PhotoHeight = 0

	err = conn.Create(&tPhoto).Error
	if err != nil {
		return tPhoto, err
	}

	// save the image to a file
	err = r.saveLocal(tPhoto.PhotoPath, file)
	if err != nil {
		return tPhoto, err
	}

	return tPhoto, err
}

func (r repository) GetTableById(conn *gorm.DB, id string, preloads ...string) (tPhoto model.Photo, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&tPhoto).Error
	return tPhoto, err
}

func (r repository) GetViewById(conn *gorm.DB, id string, preloads ...string) (vPhoto model.PhotoView, err error) {
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ? ", id).First(&vPhoto).Error
	return vPhoto, err
}

func (r repository) Create(conn *gorm.DB, tPhoto model.Photo) error {
	return conn.Create(&tPhoto).Error
}

func (r repository) Update(conn *gorm.DB, tPhoto model.Photo) error {
	return conn.Model(&tPhoto).Updates(&tPhoto).Error
}

func (r repository) Save(conn *gorm.DB, tPhoto model.Photo) error {
	return conn.Save(&tPhoto).Error
}

func (r repository) Delete(conn *gorm.DB, tPhoto model.Photo) error {
	return conn.Delete(&tPhoto).Error
}

// everytime the func called add running + 1
func (r repository) photoincGettouse(conn *gorm.DB, photoRef model.PhotoRef) (tPhotoinc model.Photoinc, err error) {
	err = conn.Where("ref_table = ?", photoRef).
		Order("folder_inc DESC").
		First(&tPhotoinc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tPhotoinc, err = r.photoincNew(conn, photoRef, 1)
		}
	} else {
		if tPhotoinc.Running >= config.PhotoincRunningLimit {
			tPhotoinc, err = r.photoincNew(conn, photoRef, tPhotoinc.FolderInc+1)
		} else {
			err = r.photoincAddrunning(conn, tPhotoinc)
		}
	}

	return tPhotoinc, err
}

func (r repository) photoincNew(conn *gorm.DB, photoRef model.PhotoRef, folderInc int64) (tPhotoinc model.Photoinc, err error) {
	tPhotoinc = model.Photoinc{
		RefTable:  string(photoRef),
		FolderInc: folderInc,
		Folder:    fmt.Sprintf("%s/%s/%s/%d", config.StorageDirectory, config.PhotoDirectory, photoRef, folderInc),
		Running:   1,
	}

	err = r.createPhotoInc(conn, tPhotoinc)
	if err != nil {
		return tPhotoinc, err
	}

	err = utils.CreateFolder(tPhotoinc.Folder, 0777)
	if err != nil {
		return tPhotoinc, err
	}

	return tPhotoinc, err
}

func (r repository) photoincAddrunning(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	tPhotoinc.Running = tPhotoinc.Running + 1
	return conn.Save(&tPhotoinc).Error
}

func (r repository) createPhotoInc(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	return conn.Create(&tPhotoinc).Error
}

func (r repository) saveLocal(filepath string, file *multipart.FileHeader) error {
	return utils.UploadImage(filepath, file)
}

func (r repository) deleteLocal(filepath string) error {
	return utils.DeleteFileLocal(filepath)
}

func NewRepository() Repository {
	return repository{}
}
