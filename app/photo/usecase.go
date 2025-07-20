package photo

import (
	"fmt"
	"github.com/jihanlugas/warehouse/app/photoinc"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/utils"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
)

type Usecase interface {
	Upload(loginUser jwt.UserLogin, conn *gorm.DB, file *multipart.FileHeader, refTable model.PhotoRef) (tPhoto model.Photo, err error)
}

type usecase struct {
	photoRepository    Repository
	photoincRepository photoinc.Repository
}

func (u usecase) Upload(loginUser jwt.UserLogin, conn *gorm.DB, file *multipart.FileHeader, refTable model.PhotoRef) (tPhoto model.Photo, err error) {
	var tPhotoinc model.Photoinc

	tPhotoinc, err = u.photoincRepository.GetTableToUse(conn, refTable)
	if err != nil {
		return tPhoto, err
	}

	err = u.photoincRepository.AddRunning(conn, tPhotoinc)
	if err != nil {
		return tPhoto, err
	}

	tPhoto.ID = utils.GetUniqueID()
	tPhoto.Ext = filepath.Ext(file.Filename)
	tPhoto.ClientName = file.Filename
	tPhoto.ServerName = fmt.Sprintf("%s%s", tPhoto.ID, tPhoto.Ext)
	tPhoto.RefTable = string(refTable)
	tPhoto.PhotoPath = fmt.Sprintf("%s/%s", tPhotoinc.Folder, tPhoto.ServerName)
	tPhoto.PhotoSize = file.Size
	tPhoto.PhotoWidth = 0
	tPhoto.PhotoHeight = 0
	tPhoto.CreateBy = loginUser.UserID
	tPhoto.UpdateBy = loginUser.UserID
	err = u.photoRepository.Create(conn, tPhoto)
	if err != nil {
		return tPhoto, err
	}

	err = utils.UploadImage(tPhoto.PhotoPath, file)
	if err != nil {
		return tPhoto, err
	}

	return tPhoto, err
}

func NewUsecase(photoRepository Repository, photoincRepository photoinc.Repository) Usecase {
	return &usecase{
		photoRepository:    photoRepository,
		photoincRepository: photoincRepository,
	}
}
