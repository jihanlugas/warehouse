package request

import (
	"time"

	"github.com/jihanlugas/warehouse/model"
)

type UpdateTransferin struct {
	ReceivedGrossQuantity float64 `json:"receivedGrossQuantity" form:"receivedGrossQuantity" query:"receivedGrossQuantity" validate:""`
	ReceivedTareQuantity  float64 `json:"receivedTareQuantity" form:"receivedTareQuantity" query:"receivedTareQuantity" validate:""`
	ReceivedNetQuantity   float64 `json:"receivedNetQuantity" form:"receivedNetQuantity" query:"receivedNetQuantity" validate:""`
}

type PageTransferin struct {
	Paging
	ProductID                  string                           `json:"productId" form:"productId" query:"productId"`
	VehicleID                  string                           `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
	Notes                      string                           `json:"notes" form:"notes" query:"notes"`
	StockmovementvehicleStatus model.StockmovementvehicleStatus `json:"stockmovementvehicleStatus" form:"stockmovementvehicleStatus" query:"stockmovementvehicleStatus"`
	StartReceivedGrossQuantity *float64                         `json:"startReceivedGrossQuantity" form:"startReceivedGrossQuantity" query:"startReceivedGrossQuantity"`
	StartReceivedTareQuantity  *float64                         `json:"startReceivedTareQuantity" form:"startReceivedTareQuantity" query:"startReceivedTareQuantity"`
	StartReceivedNetQuantity   *float64                         `json:"startReceivedNetQuantity" form:"startReceivedNetQuantity" query:"startReceivedNetQuantity"`
	StartReceivedTime          *time.Time                       `json:"startReceivedTime" form:"startReceivedTime" query:"startReceivedTime"`
	EndReceivedGrossQuantity   *float64                         `json:"endReceivedGrossQuantity" form:"endReceivedGrossQuantity" query:"endReceivedGrossQuantity"`
	EndReceivedTareQuantity    *float64                         `json:"endReceivedTareQuantity" form:"endReceivedTareQuantity" query:"endReceivedTareQuantity"`
	EndReceivedNetQuantity     *float64                         `json:"endReceivedNetQuantity" form:"endReceivedNetQuantity" query:"endReceivedNetQuantity"`
	EndReceivedTime            *time.Time                       `json:"endReceivedTime" form:"endReceivedTime" query:"endReceivedTime"`
	CreateName                 string                           `json:"createName" form:"createName" query:"createName"`
	StartCreateDt              *time.Time                       `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt                *time.Time                       `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads                   string                           `json:"preloads" form:"preloads" query:"preloads"`
}
