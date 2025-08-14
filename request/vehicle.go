package request

type CreateVehicle struct {
	WarehouseID string `json:"warehouseId" form:"warehouseId" query:"warehouseId" validate:"required"`
	PlateNumber string `json:"plateNumber" form:"plateNumber" query:"plateNumber" validate:"required"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Notes       string `json:"notes" form:"notes" query:"notes" validate:""`
	NIK         string `json:"nik" form:"nik" query:"nik" validate:""`
	DriverName  string `json:"driverName" form:"driverName" query:"driverName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber" validate:"required"`
}
type UpdateVehicle struct {
	PlateNumber string `json:"plateNumber" form:"plateNumber" query:"plateNumber" validate:"required"`
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Notes       string `json:"notes" form:"notes" query:"notes" validate:""`
	NIK         string `json:"nik" form:"nik" query:"nik" validate:""`
	DriverName  string `json:"driverName" form:"driverName" query:"driverName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber" validate:"required"`
}

type PageVehicle struct {
	Paging
	WarehouseID string `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	PlateNumber string `json:"plateNumber" form:"plateNumber" query:"plateNumber"`
	Name        string `json:"name" form:"name" query:"name"`
	Notes       string `json:"notes" form:"notes" query:"notes"`
	NIK         string `json:"nik" form:"nik" query:"nik"`
	DriverName  string `json:"driverName" form:"driverName" query:"driverName"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
	CreateName  string `json:"createName" form:"createName" query:"createName"`
	Preloads    string `json:"preloads" form:"preloads" query:"preloads"`
}
