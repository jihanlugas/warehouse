package request

import "time"

type CreateStockmovementvehicle struct {
	StockmovementID      string     `json:"stockmovementId" form:"stockmovementId" query:"stockmovementId" validate:"required"`
	ProductID            string     `json:"productId" form:"productId" query:"productId" validate:"required"`
	VehicleID            string     `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:"required"`
	SentGrossQuantity    float64    `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:"required"`
	SentTareQuantity     float64    `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:"required"`
	SentNetQuantity      float64    `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:"required"`
	SentTime             *time.Time `json:"sentTime" form:"sentTime" query:"sentTime" validate:""`
	RecivedGrossQuantity float64    `json:"recivedGrossQuantity" form:"recivedGrossQuantity" query:"recivedGrossQuantity" validate:""`
	RecivedTareQuantity  float64    `json:"recivedTareQuantity" form:"recivedTareQuantity" query:"recivedTareQuantity" validate:""`
	RecivedNetQuantity   float64    `json:"recivedNetQuantity" form:"recivedNetQuantity" query:"recivedNetQuantity" validate:""`
	RecivedTime          *time.Time `json:"recivedTime" form:"recivedTime" query:"recivedTime" validate:""`
}

type UpdateStockmovementvehicle struct {
	StockmovementID      string     `json:"stockmovementId" form:"stockmovementId" query:"stockmovementId" validate:"required"`
	ProductID            string     `json:"productId" form:"productId" query:"productId" validate:"required"`
	VehicleID            string     `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:"required"`
	SentGrossQuantity    float64    `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:"required"`
	SentTareQuantity     float64    `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:"required"`
	SentNetQuantity      float64    `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:"required"`
	SentTime             *time.Time `json:"sentTime" form:"sentTime" query:"sentTime" validate:""`
	RecivedGrossQuantity float64    `json:"recivedGrossQuantity" form:"recivedGrossQuantity" query:"recivedGrossQuantity" validate:""`
	RecivedTareQuantity  float64    `json:"recivedTareQuantity" form:"recivedTareQuantity" query:"recivedTareQuantity" validate:""`
	RecivedNetQuantity   float64    `json:"recivedNetQuantity" form:"recivedNetQuantity" query:"recivedNetQuantity" validate:""`
	RecivedTime          *time.Time `json:"recivedTime" form:"recivedTime" query:"recivedTime" validate:""`
}

type PageStockmovementvehicle struct {
	Paging
	StockmovementID           string     `json:"stockmovementId" form:"stockmovementId" query:"stockmovementId"`
	ProductID                 string     `json:"productId" form:"productId" query:"productId"`
	VehicleID                 string     `json:"vehicleId" form:"vehicleId" query:"vehicleId"`
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
