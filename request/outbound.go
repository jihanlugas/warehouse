package request

import "time"

type CreateOutbound struct {
	IsDirect               bool       `json:"isDirect" from:"isDirect" query:"isDirect" validate:""`
	IsNewVehiclerdriver    bool       `json:"isNewVehiclerdriver" from:"isNewVehiclerdriver" query:"isNewVehiclerdriver" validate:""`
	PlateNumber            string     `json:"plateNumber" from:"plateNumber" query:"plateNumber" validate:""`
	VehicleName            string     `json:"vehicleName" from:"vehicleName" query:"vehicleName" validate:""`
	NIK                    string     `json:"nik" from:"nik" query:"nik" validate:""`
	DriverName             string     `json:"driverName" from:"driverName" query:"driverName" validate:""`
	PhoneNumber            string     `json:"phoneNumber" from:"phoneNumber" query:"phoneNumber" validate:""`
	FromWarehouseID        string     `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId" validate:"required"`
	ToWarehouseID          string     `json:"toWarehouseId" form:"toWarehouseId" query:"toWarehouseId" validate:"required"`
	Remark                 string     `json:"remark" form:"remark" query:"remark" validate:""`
	ProductID              string     `json:"productId" form:"productId" query:"productId" validate:"required"`
	VehicleID              string     `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:""`
	StockmovementvehicleID string     `json:"stockmovementvehicleId" form:"stockmovementvehicleId" query:"stockmovementvehicleId" validate:""`
	SentGrossQuantity      float64    `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity       float64    `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity        float64    `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
	SentTime               *time.Time `json:"sentTime" form:"sentTime" query:"sentTime" validate:""`
}

type UpdateOutbound struct {
	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}
type PageOutbound struct {
	Paging
	WarehouseID               string     `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	StockmovementID           string     `json:"stockmovementId" form:"stockmovementId" query:"stockmovementId"`
	ProductID                 string     `json:"productId" form:"productId" query:"productId"`
	VehicleID                 string     `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
	Type                      string     `json:"type" form:"type" query:"type"`
	Remark                    string     `json:"remark" form:"remark" query:"remark"`
	Status                    string     `json:"status" form:"status" query:"status"`
	StartSentGrossQuantity    *float64   `json:"startSentGrossQuantity" form:"startSentGrossQuantity" query:"startSentGrossQuantity"`
	StartSentTareQuantity     *float64   `json:"startSentTareQuantity" form:"startSentTareQuantity" query:"startSentTareQuantity"`
	StartSentNetQuantity      *float64   `json:"startSentNetQuantity" form:"startSentNetQuantity" query:"startSentNetQuantity"`
	StartSentTime             *time.Time `json:"startSentTime" form:"startSentTime" query:"startSentTime"`
	StartRecivedGrossQuantity *float64   `json:"startRecivedGrossQuantity" form:"startRecivedGrossQuantity" query:"startRecivedGrossQuantity"`
	StartRecivedTareQuantity  *float64   `json:"startRecivedTareQuantity" form:"startRecivedTareQuantity" query:"startRecivedTareQuantity"`
	StartRecivedNetQuantity   *float64   `json:"startRecivedNetQuantity" form:"startRecivedNetQuantity" query:"startRecivedNetQuantity"`
	StartRecivedTime          *time.Time `json:"startRecivedTime" form:"startRecivedTime" query:"startRecivedTime"`
	EndSentGrossQuantity      *float64   `json:"endSentGrossQuantity" form:"endSentGrossQuantity" query:"endSentGrossQuantity"`
	EndSentTareQuantity       *float64   `json:"endSentTareQuantity" form:"endSentTareQuantity" query:"endSentTareQuantity"`
	EndSentNetQuantity        *float64   `json:"endSentNetQuantity" form:"endSentNetQuantity" query:"endSentNetQuantity"`
	EndSentTime               *time.Time `json:"endSentTime" form:"endSentTime" query:"endSentTime"`
	EndRecivedGrossQuantity   *float64   `json:"endRecivedGrossQuantity" form:"endRecivedGrossQuantity" query:"endRecivedGrossQuantity"`
	EndRecivedTareQuantity    *float64   `json:"endRecivedTareQuantity" form:"endRecivedTareQuantity" query:"endRecivedTareQuantity"`
	EndRecivedNetQuantity     *float64   `json:"endRecivedNetQuantity" form:"endRecivedNetQuantity" query:"endRecivedNetQuantity"`
	EndRecivedTime            *time.Time `json:"endRecivedTime" form:"endRecivedTime" query:"endRecivedTime"`
	CreateName                string     `json:"createName" form:"createName" query:"createName"`
	Preloads                  string     `json:"preloads" form:"preloads" query:"preloads"`
}
