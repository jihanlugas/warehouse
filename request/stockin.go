package request

import "time"

type CreateStockin struct {
	WarehouseID   string  `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	ProductID     string  `json:"productId" form:"productId" query:"productId"`
	Remark        string  `json:"remark" form:"remark" query:"remark"`
	GrossQuantity float64 `json:"grossQuantity" form:"grossQuantity" query:"grossQuantity" validate:""`
	TareQuantity  float64 `json:"tareQuantity" form:"tareQuantity" query:"tareQuantity" validate:""`
	NetQuantity   float64 `json:"netQuantity" form:"netQuantity" query:"netQuantity" validate:"required"`
}

type PageStockin struct {
	Paging
	WarehouseID        string     `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	ProductID          string     `json:"productId" form:"productId" query:"productId"`
	Remark             string     `json:"remark" form:"remark" query:"remark"`
	StartGrossQuantity *float64   `json:"startGrossQuantity" form:"startGrossQuantity" query:"startGrossQuantity"`
	StartTareQuantity  *float64   `json:"startTareQuantity" form:"startTareQuantity" query:"startTareQuantity"`
	StartNetQuantity   *float64   `json:"startNetQuantity" form:"startNetQuantity" query:"startNetQuantity"`
	StartCreateDt      *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndGrossQuantity   *float64   `json:"endGrossQuantity" form:"endGrossQuantity" query:"endGrossQuantity"`
	EndTareQuantity    *float64   `json:"endTareQuantity" form:"endTareQuantity" query:"endTareQuantity"`
	EndNetQuantity     *float64   `json:"endNetQuantity" form:"endNetQuantity" query:"endNetQuantity"`
	EndCreateDt        *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	CreateName         string     `json:"createName" form:"createName" query:"createName"`
	Preloads           string     `json:"preloads" form:"preloads" query:"preloads"`
}
