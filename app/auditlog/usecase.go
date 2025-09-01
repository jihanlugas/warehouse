package auditlog

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jihanlugas/warehouse/db"
	"github.com/jihanlugas/warehouse/jwt"
	"github.com/jihanlugas/warehouse/model"
	"github.com/jihanlugas/warehouse/request"
	"github.com/jihanlugas/warehouse/utils"
)

type Usecase interface {
	Page(loginUser jwt.UserLogin, req request.PageAuditlog) (vAuditlogs []model.AuditlogView, count int64, err error)
	GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vAuditlog model.AuditlogView, err error)
	CreateAuditlog(loginUser jwt.UserLogin, auditlogtype model.AuditlogType, req request.CreateAuditlog) error
}

type usecase struct {
	auditlogRepository Repository
}

func (u usecase) Page(loginUser jwt.UserLogin, req request.PageAuditlog) (vAuditlogs []model.AuditlogView, count int64, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vAuditlogs, count, err = u.auditlogRepository.Page(conn, req)
	if err != nil {
		return vAuditlogs, count, err
	}

	return vAuditlogs, count, err
}

func (u usecase) GetById(loginUser jwt.UserLogin, id string, preloads ...string) (vAuditlog model.AuditlogView, err error) {
	conn, closeConn := db.GetConnection()
	defer closeConn()

	vAuditlog, err = u.auditlogRepository.GetViewById(conn, id, preloads...)
	if err != nil {
		return vAuditlog, errors.New(fmt.Sprintf("failed to get %s: %v", u.auditlogRepository.Name(), err))
	}

	return vAuditlog, err
}

func (u usecase) CreateAuditlog(loginUser jwt.UserLogin, auditlogtype model.AuditlogType, req request.CreateAuditlog) error {
	var err error
	var tAuditlog model.Auditlog

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	bytesRequest, _ := json.Marshal(req.Request)
	bytesResponse, _ := json.Marshal(req.Response)

	tAuditlog = model.Auditlog{
		ID:                     utils.GetUniqueID(),
		StockmovementvehicleID: req.StockmovementvehicleID,
		LocationID:             loginUser.LocationID,
		WarehouseID:            loginUser.WarehouseID,
		AuditlogType:           auditlogtype,
		Title:                  req.Title,
		Description:            req.Description,
		Request:                string(bytesRequest),
		Response:               string(bytesResponse),
		CreateBy:               loginUser.UserID,
		UpdateBy:               loginUser.UserID,
	}

	err = u.auditlogRepository.Create(tx, tAuditlog)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create %s: %v", u.auditlogRepository.Name(), err))
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func NewUsecase(auditlogRepository Repository) Usecase {
	return &usecase{
		auditlogRepository: auditlogRepository,
	}
}
