package request

import "time"

type ChangePassword struct {
	CurrentPasswd string `json:"currentPasswd" form:"currentPasswd" query:"currentPasswd" validate:"required,lte=200"`
	Passwd        string `json:"passwd" form:"passwd" query:"passwd" validate:"required,lte=200"`
	ConfirmPasswd string `json:"confirmPasswd" form:"confirmPasswd" query:"confirmPasswd" validate:"required,lte=200,eqfield=Passwd"`
}

type CreateUser struct {
	WarehouseID   string     `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	Fullname      string     `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email         string     `json:"email" form:"email" validate:"required,lte=200,email,notexists=email"`
	UserRole      string     `json:"userRole" form:"userRole" validate:"required"`
	PhoneNumber   string     `json:"phoneNumber" form:"phoneNumber" validate:"required,lte=20"`
	Username      string     `json:"username" form:"username" validate:"required,lte=20,lowercase,notexists=username"`
	Passwd        string     `json:"passwd" form:"passwd" validate:"required,lte=200"`
	Address       string     `json:"address" form:"address" validate:""`
	BirthDt       *time.Time `json:"birthDt" form:"birthDt" validate:""`
	BirthPlace    string     `json:"birthPlace" form:"birthPlace" validate:""`
	StockIn       bool       `json:"stockIn" form:"stockIn" query:"stockIn"`
	TransferOut   bool       `json:"transferOut" form:"transferOut" query:"transferOut"`
	TransferIn    bool       `json:"transferIn" form:"transferIn" query:"transferIn"`
	Purchaseorder bool       `json:"purchaseorder" form:"purchaseorder" query:"purchaseorder"`
	Retail        bool       `json:"retail" form:"retail" query:"retail"`
}

type UpdateUser struct {
	Fullname    string     `json:"fullname" form:"fullname" validate:"required,lte=80"`
	Email       string     `json:"email" form:"email" validate:"required,lte=200,email"`
	PhoneNumber string     `json:"phoneNumber" form:"phoneNumber" validate:"required,lte=20"`
	Username    string     `json:"username" form:"username" validate:"required,lte=20,lowercase"`
	Address     string     `json:"address" form:"address" validate:""`
	BirthDt     *time.Time `json:"birthDt" form:"birthDt" validate:""`
	BirthPlace  string     `json:"birthPlace" form:"birthPlace" validate:""`
}

type PageUser struct {
	Paging
	WarehouseID   string     `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	UserRole      string     `json:"userRole" form:"userRole" query:"userRole"`
	Fullname      string     `json:"fullname" form:"fullname" query:"fullname"`
	Email         string     `json:"email" form:"email" query:"email"`
	PhoneNumber   string     `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
	Username      string     `json:"username" form:"username" query:"username"`
	Address       string     `json:"address" form:"address" query:"address"`
	BirthPlace    string     `json:"birthPlace" form:"birthPlace" query:"birthPlace"`
	CreateName    string     `json:"createName" form:"createName" query:"createName"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
