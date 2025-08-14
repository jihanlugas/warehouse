package request

import (
	"time"

	"github.com/jihanlugas/warehouse/model"
)

type CreateStockin struct {
	ProductID   string  `json:"productId" form:"productId" query:"productId" validate:"required"`
	Notes       string  `json:"notes" form:"notes" query:"notes"`
	NetQuantity float64 `json:"netQuantity" form:"netQuantity" query:"netQuantity" validate:"required"`
}

type PageStockin struct {
	Paging
	ProductID                  string                           `json:"productId" form:"productId" query:"productId"`
	StockmovementvehicleStatus model.StockmovementvehicleStatus `json:"stockmovementvehicleStatus" form:"stockmovementvehicleStatus" query:"stockmovementvehicleStatus"`
	Notes                      string                           `json:"notes" form:"notes" query:"notes"`
	StartNetQuantity           *float64                         `json:"startNetQuantity" form:"startNetQuantity" query:"startNetQuantity"`
	StartCreateDt              *time.Time                       `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndNetQuantity             *float64                         `json:"endNetQuantity" form:"endNetQuantity" query:"endNetQuantity"`
	EndCreateDt                *time.Time                       `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	CreateName                 string                           `json:"createName" form:"createName" query:"createName"`
	Preloads                   string                           `json:"preloads" form:"preloads" query:"preloads"`
}
