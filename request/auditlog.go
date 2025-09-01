package request

import "time"

type CreateAuditlog struct {
	StockmovementvehicleID string      `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId" validate:""`
	Title                  string      `json:"title" form:"title" query:"title" validate:""`
	Description            string      `json:"description" form:"description" query:"description" validate:""`
	Request                interface{} `json:"request" form:"request" query:"request" validate:""`
	Response               interface{} `json:"response" form:"response" query:"response" validate:""`
}

type PageAuditlog struct {
	Paging
	LocationID    string     `json:"locationId" form:"locationId" query:"locationId"`
	WarehouseID   string     `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	AuditlogType  string     `json:"auditlogType" form:"auditlogType" query:"auditlogType"`
	Title         string     `json:"title" form:"title" query:"title"`
	Description   string     `json:"description" form:"description" query:"description"`
	CreateName    string     `json:"createName" form:"createName" query:"createName"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
