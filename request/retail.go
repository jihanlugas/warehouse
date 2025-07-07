package request

import "time"

type CreateRetail struct {
	IsNewCustomer       bool            `json:"isNewCustomer" form:"isNewCustomer" query:"isNewCustomer" validate:""`
	CustomerID          string          `json:"customerId" form:"customerId" query:"customerId" validate:""`
	CustomerName        string          `json:"customerName" form:"customerName" query:"customerName" validate:""`
	CustomerPhoneNumber string          `json:"customerPhoneNumber" form:"customerPhoneNumber" query:"customerPhoneNumber" validate:""`
	Notes               string          `json:"notes" form:"notes" query:"notes" validate:""`
	Products            []RetailProduct `json:"products" form:"products" query:"products" validate:""`
}

type RetailProduct struct {
	ProductID string  `json:"productID" form:"productID" query:"productID" validate:"required"`
	UnitPrice float64 `json:"unitPrice" form:"" query:"" validate:"required"`
}

type UpdateRetail struct {
	Notes string `json:"notes" form:"notes" query:"notes" validate:""`
}

type CreateRetailStockmovementvehicle struct {
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

type UpdateRetailStockmomentvehicle struct {
	SentGrossQuantity float64 `json:"sentGrossQuantity" form:"sentGrossQuantity" query:"sentGrossQuantity" validate:""`
	SentTareQuantity  float64 `json:"sentTareQuantity" form:"sentTareQuantity" query:"sentTareQuantity" validate:""`
	SentNetQuantity   float64 `json:"sentNetQuantity" form:"sentNetQuantity" query:"sentNetQuantity" validate:""`
}

type PageRetail struct {
	Paging
	CustomerID        string     `json:"customeId" form:"customeId" query:"customeId"`
	Notes             string     `json:"notes" form:"notes" query:"notes"`
	Number            string     `json:"number" form:"number" query:"number"`
	Status            *string    `json:"status" form:"status" query:"status"`
	CreateName        string     `json:"createName" form:"createName" query:"createName"`
	StartTotalPrice   *float64   `json:"startTotalPrice" form:"startTotalPrice" query:"startTotalPrice"`
	EndTotalPrice     *float64   `json:"endTotalPrice" form:"endTotalPrice" query:"endTotalPrice"`
	StartTotalPayment *float64   `json:"startTotalPayment" form:"startTotalPayment" query:"startTotalPayment"`
	EndTotalPayment   *float64   `json:"endTotalPayment" form:"endTotalPayment" query:"endTotalPayment"`
	StartOutstanding  *float64   `json:"startOutstanding" form:"startOutstanding" query:"startOutstanding"`
	EndOutstanding    *float64   `json:"endOutstanding" form:"endOutstanding" query:"endOutstanding"`
	StartCreateDt     *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt       *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads          string     `json:"preloads" form:"preloads" query:"preloads"`
}
