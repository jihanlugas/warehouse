package request

import (
	"time"

	"github.com/jihanlugas/warehouse/model"
)

type CreateStockmovementvehicleRetail struct {
	IsNewVehiclerdriver bool    `json:"isNewVehiclerdriver" from:"isNewVehiclerdriver" query:"isNewVehiclerdriver" validate:""`
	PlateNumber         string  `json:"plateNumber" from:"plateNumber" query:"plateNumber" validate:""`
	VehicleName         string  `json:"vehicleName" from:"vehicleName" query:"vehicleName" validate:""`
	NIK                 string  `json:"nik" from:"nik" query:"nik" validate:""`
	DriverName          string  `json:"driverName" from:"driverName" query:"driverName" validate:""`
	PhoneNumber         string  `json:"phoneNumber" from:"phoneNumber" query:"phoneNumber" validate:""`
	RetailID            string  `json:"retailId" form:"retailId" query:"retailId" validate:"required"`
	Notes               string  `json:"notes" form:"notes" query:"notes"`
	ProductID           string  `json:"productId" form:"productId" query:"productId" validate:"required"`
	VehicleID           string  `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:""`
	SentGrossQuantity   float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity    float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity     float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}

type UpdateStockmovementvehicleRetail struct {
	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}

type PageStockmovementvehicleRetail struct {
	Paging
	RetailID                   string                           `json:"retailId" form:"retailId" query:"retailId"`
	ProductID                  string                           `json:"productId" form:"productId" query:"productId"`
	VehicleID                  string                           `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
	Notes                      string                           `json:"notes" form:"notes" query:"notes"`
	StockmovementvehicleStatus model.StockmovementvehicleStatus `json:"stockmovementvehicleStatus" form:"stockmovementvehicleStatus" query:"stockmovementvehicleStatus"`
	StartSentGrossQuantity     *float64                         `json:"startSentGrossQuantity" form:"startSentGrossQuantity" query:"startSentGrossQuantity"`
	StartSentTareQuantity      *float64                         `json:"startSentTareQuantity" form:"startSentTareQuantity" query:"startSentTareQuantity"`
	StartSentNetQuantity       *float64                         `json:"startSentNetQuantity" form:"startSentNetQuantity" query:"startSentNetQuantity"`
	StartSentTime              *time.Time                       `json:"startSentTime" form:"startSentTime" query:"startSentTime"`
	EndSentGrossQuantity       *float64                         `json:"endSentGrossQuantity" form:"endSentGrossQuantity" query:"endSentGrossQuantity"`
	EndSentTareQuantity        *float64                         `json:"endSentTareQuantity" form:"endSentTareQuantity" query:"endSentTareQuantity"`
	EndSentNetQuantity         *float64                         `json:"endSentNetQuantity" form:"endSentNetQuantity" query:"endSentNetQuantity"`
	EndSentTime                *time.Time                       `json:"endSentTime" form:"endSentTime" query:"endSentTime"`
	CreateName                 string                           `json:"createName" form:"createName" query:"createName"`
	StartCreateDt              *time.Time                       `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt                *time.Time                       `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads                   string                           `json:"preloads" form:"preloads" query:"preloads"`
}
