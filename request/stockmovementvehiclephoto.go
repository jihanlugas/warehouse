package request

import "mime/multipart"

type CreateStockmovementvehiclephoto struct {
	WarehouseID            string                `json:"warehouseId" form:"warehouseId" query:"warehouseId" validate:"required"`
	StockmovementvehicleID string                `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId" validate:"required"`
	Photo                  *multipart.FileHeader `json:"-"`
	PhotoChk               bool                  `json:"photo" form:"photo" query:"photo" validate:"required,photo=Photo"`
}

type UpdateStockmovementvehiclephoto struct {
}

type PageStockmovementvehiclephoto struct {
	Paging
	WarehouseID            string `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	StockmovementvehicleID string `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId"`
	CreateName             string `json:"createName" form:"createName" query:"createName"`
	Preloads               string `json:"preloads" form:"preloads" query:"preloads"`
}
