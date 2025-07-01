package request

import "time"

type CreatePurchaseorder struct {
	IsNewCustomer       bool                   `json:"isNewCustomer" form:"isNewCustomer" query:"isNewCustomer" validate:""`
	CustomerID          string                 `json:"customerId" form:"customerId" query:"customerId" validate:""`
	CustomerName        string                 `json:"customerName" form:"customerName" query:"customerName" validate:""`
	CustomerPhoneNumber string                 `json:"customerPhoneNumber" form:"customerPhoneNumber" query:"customerPhoneNumber" validate:""`
	Notes               string                 `json:"notes" form:"notes" query:"notes" validate:""`
	TotalAmount         float64                `json:"" form:"" query:"" validate:""`
	Products            []PurchaseorderProduct `json:"products" form:"products" query:"products" validate:""`
}

type PurchaseorderProduct struct {
	ProductID string  `json:"productID" form:"productID" query:"productID" validate:"required"`
	UnitPrice float64 `json:"unitPrice" form:"" query:"" validate:"required"`
}

type UpdatePurchaseorder struct {
	Notes string `json:"notes" form:"notes" query:"notes" validate:""`
}

type CreatePurchaseorderStockmovementvehicle struct {
	IsNewVehiclerdriver bool    `json:"isNewVehiclerdriver" from:"isNewVehiclerdriver" query:"isNewVehiclerdriver" validate:""`
	PlateNumber         string  `json:"plateNumber" from:"plateNumber" query:"plateNumber" validate:""`
	VehicleName         string  `json:"vehicleName" from:"vehicleName" query:"vehicleName" validate:""`
	NIK                 string  `json:"nik" from:"nik" query:"nik" validate:""`
	DriverName          string  `json:"driverName" from:"driverName" query:"driverName" validate:""`
	PhoneNumber         string  `json:"phoneNumber" from:"phoneNumber" query:"phoneNumber" validate:""`
	VehicleID           string  `json:"vehicleId" form:"vehicleId" query:"vehicleId" validate:"required"`
	StockmovementID     string  `json:"stockmovementId" form:"stockmovementId" query:"stockmovementId" validate:"required"`
	ProductID           string  `json:"productId" form:"productId" query:"productId" validate:"required"`
	SentGrossQuantity   float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:"required"`
	SentTareQuantity    float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:"required"`
	SentNetQuantity     float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:"required"`
}

type UpdatePurchaseorderStockmomentvehicle struct {
	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}

type PagePurchaseorder struct {
	Paging
	CustomerID       string     `json:"customeId" form:"customeId" query:"customeId"`
	Notes            string     `json:"notes" form:"notes" query:"notes"`
	CreateName       string     `json:"createName" form:"createName" query:"createName"`
	StartTotalAmount *float64   `json:"startTotalAmount" form:"startTotalAmount" query:"startTotalAmount"`
	EndTotalAmount   *float64   `json:"endTotalAmount" form:"endTotalAmount" query:"endTotalAmount"`
	StartCreateDt    *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt      *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads         string     `json:"preloads" form:"preloads" query:"preloads"`
}
