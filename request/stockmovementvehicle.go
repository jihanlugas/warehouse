package request

import (
	"time"

	"github.com/jihanlugas/warehouse/model"
)

//type CreateStockmovementvehiclePurchaseorder struct {
//	IsDirect               bool    `json:"isDirect" from:"isDirect" query:"isDirect" validate:""`
//	FromWarehouseID        string  `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId" validate:"required"`
//	PurchaseorderID        string  `json:"purchaseorderId" form:"purchaseorderId" query:"purchaseorderId" validate:"required"`
//	ProductID              string  `json:"productId" form:"productId" query:"productId" validate:"required"`
//	StockmovementvehicleID string  `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId" validate:""`
//	IsNewVehiclerdriver    bool    `json:"isNewVehiclerdriver" from:"isNewVehiclerdriver" query:"isNewVehiclerdriver" validate:""`
//	PlateNumber            string  `json:"plateNumber" from:"plateNumber" query:"plateNumber" validate:""`
//	VehicleID              string  `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:""`
//	VehicleName            string  `json:"vehicleName" from:"vehicleName" query:"vehicleName" validate:""`
//	NIK                    string  `json:"nik" from:"nik" query:"nik" validate:""`
//	DriverName             string  `json:"driverName" from:"driverName" query:"driverName" validate:""`
//	PhoneNumber            string  `json:"phoneNumber" from:"phoneNumber" query:"phoneNumber" validate:""`
//	SentGrossQuantity      float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
//	SentTareQuantity       float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
//	SentNetQuantity        float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
//}
//
//type UpdateStockmovementvehiclePurchaseorder struct {
//	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
//	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
//	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
//}

//type CreateStockmovementvehicleRetail struct {
//	IsDirect               bool    `json:"isDirect" from:"isDirect" query:"isDirect" validate:""`
//	FromWarehouseID        string  `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId" validate:"required"`
//	RetailID               string  `json:"retailId" form:"retailId" query:"retailId" validate:"required"`
//	ProductID              string  `json:"productId" form:"productId" query:"productId" validate:"required"`
//	StockmovementvehicleID string  `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId" validate:""`
//	IsNewVehiclerdriver    bool    `json:"isNewVehiclerdriver" from:"isNewVehiclerdriver" query:"isNewVehiclerdriver" validate:""`
//	PlateNumber            string  `json:"plateNumber" from:"plateNumber" query:"plateNumber" validate:""`
//	VehicleID              string  `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:""`
//	VehicleName            string  `json:"vehicleName" from:"vehicleName" query:"vehicleName" validate:""`
//	NIK                    string  `json:"nik" from:"nik" query:"nik" validate:""`
//	DriverName             string  `json:"driverName" from:"driverName" query:"driverName" validate:""`
//	PhoneNumber            string  `json:"phoneNumber" from:"phoneNumber" query:"phoneNumber" validate:""`
//	SentGrossQuantity      float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
//	SentTareQuantity       float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
//	SentNetQuantity        float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
//}
//
//type UpdateStockmovementvehicleRetail struct {
//	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
//	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
//	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
//}

type PageStockmovementvehicle struct {
	Paging
	FromWarehouseID            string                           `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId"`
	ToWarehouseID              string                           `json:"toWarehouseId" form:"toWarehouseId" query:"toWarehouseId"`
	ProductID                  string                           `json:"productId" form:"productId" query:"productId"`
	VehicleID                  string                           `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
	RelatedID                  string                           `json:"relatedId" form:"relatedId" query:"relatedId"`
	Number                     string                           `json:"number" form:"number" query:"number"`
	StockmovementvehicleType   model.StockmovementvehicleType   `json:"stockmovementvehicleType" form:"stockmovementvehicleType" query:"stockmovementvehicleType"`
	Notes                      string                           `json:"notes" form:"notes" query:"notes"`
	StockmovementvehicleStatus model.StockmovementvehicleStatus `json:"stockmovementvehicleStatus" form:"stockmovementvehicleStatus" query:"stockmovementvehicleStatus"`
	StartSentGrossQuantity     *float64                         `json:"startSentGrossQuantity" form:"startSentGrossQuantity" query:"startSentGrossQuantity"`
	StartSentTareQuantity      *float64                         `json:"startSentTareQuantity" form:"startSentTareQuantity" query:"startSentTareQuantity"`
	StartSentNetQuantity       *float64                         `json:"startSentNetQuantity" form:"startSentNetQuantity" query:"startSentNetQuantity"`
	StartSentTime              *time.Time                       `json:"startSentTime" form:"startSentTime" query:"startSentTime"`
	StartReceivedGrossQuantity *float64                         `json:"startReceivedGrossQuantity" form:"startReceivedGrossQuantity" query:"startReceivedGrossQuantity"`
	StartReceivedTareQuantity  *float64                         `json:"startReceivedTareQuantity" form:"startReceivedTareQuantity" query:"startReceivedTareQuantity"`
	StartReceivedNetQuantity   *float64                         `json:"startReceivedNetQuantity" form:"startReceivedNetQuantity" query:"startReceivedNetQuantity"`
	StartReceivedTime          *time.Time                       `json:"startReceivedTime" form:"startReceivedTime" query:"startReceivedTime"`
	EndSentGrossQuantity       *float64                         `json:"endSentGrossQuantity" form:"endSentGrossQuantity" query:"endSentGrossQuantity"`
	EndSentTareQuantity        *float64                         `json:"endSentTareQuantity" form:"endSentTareQuantity" query:"endSentTareQuantity"`
	EndSentNetQuantity         *float64                         `json:"endSentNetQuantity" form:"endSentNetQuantity" query:"endSentNetQuantity"`
	EndSentTime                *time.Time                       `json:"endSentTime" form:"endSentTime" query:"endSentTime"`
	EndReceivedGrossQuantity   *float64                         `json:"endReceivedGrossQuantity" form:"endReceivedGrossQuantity" query:"endReceivedGrossQuantity"`
	EndReceivedTareQuantity    *float64                         `json:"endReceivedTareQuantity" form:"endReceivedTareQuantity" query:"endReceivedTareQuantity"`
	EndReceivedNetQuantity     *float64                         `json:"endReceivedNetQuantity" form:"endReceivedNetQuantity" query:"endReceivedNetQuantity"`
	EndReceivedTime            *time.Time                       `json:"endReceivedTime" form:"endReceivedTime" query:"endReceivedTime"`
	CreateName                 string                           `json:"createName" form:"createName" query:"createName"`
	StartCreateDt              *time.Time                       `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt                *time.Time                       `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads                   string                           `json:"preloads" form:"preloads" query:"preloads"`
}
