package request

import (
	"github.com/jihanlugas/warehouse/model"
	"time"
)

type CreateStockmovementvehiclePurchaseorder struct {
	FromWarehouseID     string  `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId" validate:"required"`
	PurchaseorderID     string  `json:"purchaseorderId" form:"purchaseorderId" query:"purchaseorderId" validate:"required"`
	ProductID           string  `json:"productId" form:"productId" query:"productId" validate:"required"`
	IsNewVehiclerdriver bool    `json:"isNewVehiclerdriver" from:"isNewVehiclerdriver" query:"isNewVehiclerdriver" validate:""`
	PlateNumber         string  `json:"plateNumber" from:"plateNumber" query:"plateNumber" validate:""`
	VehicleID           string  `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:""`
	VehicleName         string  `json:"vehicleName" from:"vehicleName" query:"vehicleName" validate:""`
	NIK                 string  `json:"nik" from:"nik" query:"nik" validate:""`
	DriverName          string  `json:"driverName" from:"driverName" query:"driverName" validate:""`
	PhoneNumber         string  `json:"phoneNumber" from:"phoneNumber" query:"phoneNumber" validate:""`
	SentGrossQuantity   float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity    float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity     float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}

type UpdateStockmovementvehiclePurchaseorder struct {
	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}

type PageStockmovementvehicle struct {
	Paging
	FromWarehouseID           string                  `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId"`
	ToWarehouseID             string                  `json:"toWarehouseId" form:"toWarehouseId" query:"toWarehouseId"`
	StockmovementID           string                  `json:"stockmovementId" form:"stockmovementId" query:"stockmovementId"`
	ProductID                 string                  `json:"productId" form:"productId" query:"productId"`
	VehicleID                 string                  `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
	RelatedID                 string                  `json:"relatedId" form:"relatedId" query:"relatedId"`
	Type                      model.StockMovementType `json:"type" form:"type" query:"type"`
	StartSentGrossQuantity    *float64                `json:"startSentGrossQuantity" form:"startSentGrossQuantity" query:"startSentGrossQuantity"`
	StartSentTareQuantity     *float64                `json:"startSentTareQuantity" form:"startSentTareQuantity" query:"startSentTareQuantity"`
	StartSentNetQuantity      *float64                `json:"startSentNetQuantity" form:"startSentNetQuantity" query:"startSentNetQuantity"`
	StartSentTime             *time.Time              `json:"startSentTime" form:"startSentTime" query:"startSentTime"`
	StartRecivedGrossQuantity *float64                `json:"startRecivedGrossQuantity" form:"startRecivedGrossQuantity" query:"startRecivedGrossQuantity"`
	StartRecivedTareQuantity  *float64                `json:"startRecivedTareQuantity" form:"startRecivedTareQuantity" query:"startRecivedTareQuantity"`
	StartRecivedNetQuantity   *float64                `json:"startRecivedNetQuantity" form:"startRecivedNetQuantity" query:"startRecivedNetQuantity"`
	StartRecivedTime          *time.Time              `json:"startRecivedTime" form:"startRecivedTime" query:"startRecivedTime"`
	EndSentGrossQuantity      *float64                `json:"endSentGrossQuantity" form:"endSentGrossQuantity" query:"endSentGrossQuantity"`
	EndSentTareQuantity       *float64                `json:"endSentTareQuantity" form:"endSentTareQuantity" query:"endSentTareQuantity"`
	EndSentNetQuantity        *float64                `json:"endSentNetQuantity" form:"endSentNetQuantity" query:"endSentNetQuantity"`
	EndSentTime               *time.Time              `json:"endSentTime" form:"endSentTime" query:"endSentTime"`
	EndRecivedGrossQuantity   *float64                `json:"endRecivedGrossQuantity" form:"endRecivedGrossQuantity" query:"endRecivedGrossQuantity"`
	EndRecivedTareQuantity    *float64                `json:"endRecivedTareQuantity" form:"endRecivedTareQuantity" query:"endRecivedTareQuantity"`
	EndRecivedNetQuantity     *float64                `json:"endRecivedNetQuantity" form:"endRecivedNetQuantity" query:"endRecivedNetQuantity"`
	EndRecivedTime            *time.Time              `json:"endRecivedTime" form:"endRecivedTime" query:"endRecivedTime"`
	CreateName                string                  `json:"createName" form:"createName" query:"createName"`
	Preloads                  string                  `json:"preloads" form:"preloads" query:"preloads"`
}
