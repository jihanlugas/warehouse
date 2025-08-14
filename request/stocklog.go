package request

import (
	"time"

	"github.com/jihanlugas/warehouse/model"
)

type PageStocklog struct {
	Paging
	WarehouseID            string             `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	StockID                string             `json:"stockId" form:"stockId" query:"stockId"`
	StockproductID         string             `json:"stockproductId" form:"stockproductId" query:"stockproductId"`
	StockmovementvehicleID string             `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId"`
	ProductID              string             `json:"productId" form:"productId" query:"productId"`
	VehicleID              string             `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
	StocklogType           model.StocklogType `json:"stocklogType" form:"stocklogType" query:"stocklogType"`
	StartGrossQuantity     *float64           `json:"startGrossQuantity" form:"startGrossQuantity" query:"startGrossQuantity"`
	StartTareQuantity      *float64           `json:"startTareQuantity" form:"startTareQuantity" query:"startTareQuantity"`
	StartNetQuantity       *float64           `json:"startNetQuantity" form:"startNetQuantity" query:"startNetQuantity"`
	StartCreateDt          *time.Time         `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndGrossQuantity       *float64           `json:"endGrossQuantity" form:"endGrossQuantity" query:"endGrossQuantity"`
	EndTareQuantity        *float64           `json:"endTareQuantity" form:"endTareQuantity" query:"endTareQuantity"`
	EndNetQuantity         *float64           `json:"endNetQuantity" form:"endNetQuantity" query:"endNetQuantity"`
	EndCreateDt            *time.Time         `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	CreateName             string             `json:"createName" form:"createName" query:"createName" `
	Preloads               string             `json:"preloads" form:"preloads" query:"preloads" `
}
