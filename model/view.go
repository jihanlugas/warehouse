package model

import (
	"gorm.io/gorm"
	"time"
)

type UserView struct {
	ID                string         `json:"id"`
	WarehouseID       string         `json:"warehouseId"`
	Role              UserRole       `json:"role"`
	Email             string         `json:"email"`
	Username          string         `json:"username"`
	PhoneNumber       string         `json:"phoneNumber"`
	Address           string         `json:"address"`
	Fullname          string         `json:"fullname"`
	Passwd            string         `json:"-"`
	PassVersion       int            `json:"passVersion"`
	IsActive          bool           `json:"isActive"`
	PhotoID           string         `json:"photoId"`
	PhotoUrl          string         `json:"photoUrl"`
	LastLoginDt       *time.Time     `json:"lastLoginDt"`
	BirthDt           *time.Time     `json:"birthDt"`
	BirthPlace        string         `json:"birthPlace"`
	AccountVerifiedDt *time.Time     `json:"accountVerifiedDt"`
	CreateBy          string         `json:"createBy"`
	CreateDt          time.Time      `json:"createDt"`
	UpdateBy          string         `json:"updateBy"`
	UpdateDt          time.Time      `json:"updateDt"`
	DeleteDt          gorm.DeletedAt `json:"deleteDt"`
	CreateName        string         `json:"createName"`
	UpdateName        string         `json:"updateName"`

	Userprivilege *UserprivilegeView `json:"userprivilege,omitempty" gorm:"foreignKey:UserID"`
	Warehouse     *WarehouseView     `json:"warehouse,omitempty"`
}

func (UserView) TableName() string {
	return VIEW_USER
}

type UserprivilegeView struct {
	ID            string         `json:"id"`
	UserID        string         `json:"userId"`
	StockIn       bool           `json:"stockIn"`
	TransferOut   bool           `json:"transferOut"`
	TransferIn    bool           `json:"transferIn"`
	PurchaseOrder bool           `json:"purchaseOrder"`
	Retail        bool           `json:"retail"`
	CreateBy      string         `json:"createBy"`
	CreateDt      time.Time      `json:"createDt"`
	UpdateBy      string         `json:"updateBy"`
	UpdateDt      time.Time      `json:"updateDt"`
	DeleteDt      gorm.DeletedAt `json:"deleteDt"`
	CreateName    string         `json:"createName"`
	UpdateName    string         `json:"updateName"`
}

func (UserprivilegeView) TableName() string {
	return VIEW_USERPRIVILEGE
}

type CustomerView struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phoneNumber"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`

	Retails        []RetailView        `json:"retail,omitempty" gorm:"foreignKey:CustomerID"`
	Purchaseorders []PurchaseorderView `json:"purchaseorders,omitempty" gorm:"foreignKey:CustomerID"`
}

func (CustomerView) TableName() string {
	return VIEW_CUSTOMER
}

type RetailView struct {
	ID           string         `json:"id"`
	CustomerID   string         `json:"customerId"`
	TotalPrice   float64        `json:"totalPrice"`
	TotalPayment float64        `json:"totalPayment"`
	Outstanding  float64        `json:"outstanding"`
	Number       string         `json:"number"`
	Notes        string         `json:"notes"`
	Status       RetailStatus   `json:"status"`
	CreateBy     string         `json:"createBy"`
	CreateDt     time.Time      `json:"createDt"`
	UpdateBy     string         `json:"updateBy"`
	UpdateDt     time.Time      `json:"updateDt"`
	DeleteDt     gorm.DeletedAt `json:"deleteDt"`
	CreateName   string         `json:"createName"`
	UpdateName   string         `json:"updateName"`

	Customer              *CustomerView              `json:"customer,omitempty"`
	Retailproducts        []RetailproductView        `json:"retailproducts,omitempty" gorm:"foreignKey:RetailID"`
	Transactions          []TransactionView          `json:"transactions,omitempty" gorm:"foreignKey:RelatedID"`
	Stockmovements        []StockmovementView        `json:"stockmovements,omitempty" gorm:"foreignKey:RelatedID"`
	Stockmovementvehicles []StockmovementvehicleView `json:"stockmovementvehicles,omitempty" gorm:"foreignKey:RelatedID"`
}

func (RetailView) TableName() string {
	return VIEW_RETAIL
}

type RetailproductView struct {
	ID         string         `json:"id"`
	RetailID   string         `json:"retailId"`
	ProductID  string         `json:"productId"`
	UnitPrice  float64        `json:"unitPrice"`
	CreateBy   string         `json:"createBy"`
	CreateDt   time.Time      `json:"createDt"`
	UpdateBy   string         `json:"updateBy"`
	UpdateDt   time.Time      `json:"updateDt"`
	DeleteDt   gorm.DeletedAt `json:"deleteDt"`
	CreateName string         `json:"createName"`
	UpdateName string         `json:"updateName"`

	Retail  *RetailView  `json:"retail,omitempty"`
	Product *ProductView `json:"product,omitempty"`
}

func (RetailproductView) TableName() string {
	return VIEW_RETAILPRODUCT
}

type PurchaseorderView struct {
	ID           string              `json:"id"`
	CustomerID   string              `json:"customerId"`
	TotalPrice   float64             `json:"totalPrice"`
	TotalPayment float64             `json:"totalPayment"`
	Outstanding  float64             `json:"outstanding"`
	Number       string              `json:"number"`
	Notes        string              `json:"notes"`
	Status       PurchaseorderStatus `json:"status"`
	CreateBy     string              `json:"createBy"`
	CreateDt     time.Time           `json:"createDt"`
	UpdateBy     string              `json:"updateBy"`
	UpdateDt     time.Time           `json:"updateDt"`
	DeleteDt     gorm.DeletedAt      `json:"deleteDt"`
	CreateName   string              `json:"createName"`
	UpdateName   string              `json:"updateName"`

	Customer              *CustomerView              `json:"customer,omitempty"`
	Purchaseorderproducts []PurchaseorderproductView `json:"purchaseorderproducts,omitempty" gorm:"foreignKey:PurchaseorderID"`
	Transactions          []TransactionView          `json:"transactions,omitempty" gorm:"foreignKey:RelatedID"`
	Stockmovements        []StockmovementView        `json:"stockmovements,omitempty" gorm:"foreignKey:RelatedID"`
	Stockmovementvehicles []StockmovementvehicleView `json:"stockmovementvehicles,omitempty" gorm:"foreignKey:RelatedID"`
}

func (PurchaseorderView) TableName() string {
	return VIEW_PURCHASEORDER
}

type PurchaseorderproductView struct {
	ID              string         `json:"id"`
	PurchaseorderID string         `json:"purchaseorderId"`
	ProductID       string         `json:"productId"`
	UnitPrice       float64        `json:"unitPrice"`
	CreateBy        string         `json:"createBy"`
	CreateDt        time.Time      `json:"createDt"`
	UpdateBy        string         `json:"updateBy"`
	UpdateDt        time.Time      `json:"updateDt"`
	DeleteDt        gorm.DeletedAt `json:"deleteDt"`
	CreateName      string         `json:"createName"`
	UpdateName      string         `json:"updateName"`

	Purchaseorder *PurchaseorderView `json:"purchaseorder,omitempty"`
	Product       *ProductView       `json:"product,omitempty"`
}

func (PurchaseorderproductView) TableName() string {
	return VIEW_PURCHASEORDERPRODUCT
}

type TransactionView struct {
	ID          string                 `json:"id"`
	RelatedID   string                 `json:"RelatedId"`
	RelatedType TransactionRelatedType `json:"relatedType"`
	Type        TransactionType        `json:"type"`
	CustomerID  string                 `json:"customerId"`
	Amount      float64                `json:"amount"`
	Notes       string                 `json:"notes"`
	Number      string                 `json:"number"`
	CreateBy    string                 `json:"createBy"`
	CreateDt    time.Time              `json:"createDt"`
	UpdateBy    string                 `json:"updateBy"`
	UpdateDt    time.Time              `json:"updateDt"`
	DeleteDt    gorm.DeletedAt         `json:"deleteDt"`
	CreateName  string                 `json:"createName"`
	UpdateName  string                 `json:"updateName"`

	Customer      *CustomerView      `json:"customer,omitempty"`
	Retail        *RetailView        `json:"retail,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Purchaseorder *PurchaseorderView `json:"purchaseorder,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
}

func (TransactionView) TableName() string {
	return VIEW_TRANSACTION
}

type WarehouseView struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	Location             string         `json:"location"`
	PhoneNumber          string         `json:"phoneNumber"`
	IsStockin            bool           `json:"isStockin"`
	IsInbound            bool           `json:"isInbound"`
	IsOutbound           bool           `json:"isOutbound"`
	IsRetail             bool           `json:"isRetail"`
	IsPurchaseorder      bool           `json:"isPurchaseorder"`
	PhotoID              string         `json:"photoId"`
	PhotoUrl             string         `json:"photoUrl"`
	CreateBy             string         `json:"createBy"`
	CreateDt             time.Time      `json:"createDt"`
	UpdateBy             string         `json:"updateBy"`
	UpdateDt             time.Time      `json:"updateDt"`
	DeleteDt             gorm.DeletedAt `json:"deleteDt"`
	TotalRunningOutbound float64        `json:"totalRunningOutbound"`
	TotalRunningInbound  float64        `json:"totalRunningInbound"`
	CreateName           string         `json:"createName"`
	UpdateName           string         `json:"updateName"`

	Stocks    []StockView    `json:"stocks,omitempty" gorm:"foreignKey:WarehouseID"`
	Stocklogs []StocklogView `json:"stocklogs,omitempty" gorm:"foreignKey:WarehouseID"`
}

func (WarehouseView) TableName() string {
	return VIEW_WAREHOUSE
}

type VehicleView struct {
	ID          string         `json:"id"`
	WarehouseID string         `json:"warehouseId"`
	PlateNumber string         `json:"plateNumber"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	NIK         string         `json:"nik"`
	DriverName  string         `json:"driverName"`
	PhoneNumber string         `json:"phoneNumber"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (VehicleView) TableName() string {
	return VIEW_VEHICLE
}

type ProductView struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (ProductView) TableName() string {
	return VIEW_PRODUCT
}

type StockView struct {
	ID          string         `json:"id"`
	WarehouseID string         `json:"warehouseId"`
	ProductID   string         `json:"productId"`
	Quantity    float64        `json:"quantity"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`

	Product   *ProductView   `json:"product,omitempty"`
	Warehouse *WarehouseView `json:"warehouse,omitempty"`
}

func (StockView) TableName() string {
	return VIEW_STOCK
}

type StocklogView struct {
	ID                     string         `json:"id"`
	WarehouseID            string         `json:"warehouseId"`
	StockID                string         `json:"stockId"`
	StockmovementID        string         `json:"stockmovementId"`
	StockmovementvehicleID string         `json:"stockmovementvehicleId"`
	ProductID              string         `json:"productId"`
	VehicleID              string         `json:"vehicleId"`
	Type                   StockLogType   `json:"type"`
	GrossQuantity          float64        `json:"grossQuantity"`
	TareQuantity           float64        `json:"tareQuantity"`
	NetQuantity            float64        `json:"netQuantity"`
	CurrentQuantity        float64        `json:"currentQuantity"`
	CreateBy               string         `json:"createBy"`
	CreateDt               time.Time      `json:"createDt"`
	UpdateBy               string         `json:"updateBy"`
	UpdateDt               time.Time      `json:"updateDt"`
	DeleteDt               gorm.DeletedAt `json:"deleteDt"`
	CreateName             string         `json:"createName"`
	UpdateName             string         `json:"updateName"`

	Warehouse            *WarehouseView            `json:"warehouse,omitempty"`
	Stock                *StockView                `json:"stock,omitempty"`
	Stockmovement        *StockmovementView        `json:"stockmovement,omitempty"`
	Stockmovementvehicle *StockmovementvehicleView `json:"stockmovementvehicle,omitempty"`
	Product              *ProductView              `json:"product,omitempty"`
	Vehicle              *VehicleView              `json:"vehicle,omitempty"`
}

func (StocklogView) TableName() string {
	return VIEW_STOCKLOG
}

type StockmovementView struct {
	ID              string            `json:"id"`
	FromWarehouseID string            `json:"fromWarehouseId"`
	ToWarehouseID   string            `json:"toWarehouseId"`
	ProductID       string            `json:"productId"`
	RelatedID       string            `json:"relatedId"`
	Type            StockMovementType `json:"type"`
	UnitPrice       float64           `json:"unitPrice"`
	Remark          string            `json:"remark"`
	CreateDt        time.Time         `json:"createDt"`
	UpdateBy        string            `json:"updateBy"`
	UpdateDt        time.Time         `json:"updateDt"`
	DeleteDt        gorm.DeletedAt    `json:"deleteDt"`
	CreateName      string            `json:"createName"`
	UpdateName      string            `json:"updateName"`

	FromWarehouse *WarehouseView     `json:"fromWarehouse,omitempty" gorm:"foreignKey:FromWarehouseID;references:ID"`
	ToWarehouse   *WarehouseView     `json:"toWarehouse,omitempty" gorm:"foreignKey:ToWarehouseID;references:ID"`
	Purchaseorder *PurchaseorderView `json:"purchaseorder,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Retail        *RetailView        `json:"retail,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Product       *ProductView       `json:"product,omitempty"`
}

func (StockmovementView) TableName() string {
	return VIEW_STOCKMOVEMENT
}

type StockmovementvehicleView struct {
	ID                   string            `json:"id"`
	StockmovementID      string            `json:"stockmovementId"`
	FromWarehouseID      string            `json:"fromWarehouseId"`
	ToWarehouseID        string            `json:"toWarehouseId"`
	RelatedID            string            `json:"relatedId"`
	Type                 StockMovementType `json:"type"`
	UnitPrice            float64           `json:"unitPrice"`
	ProductID            string            `json:"productId"`
	VehicleID            string            `json:"vehicleId"`
	SentGrossQuantity    float64           `json:"sentGrossQuantity"`
	SentTareQuantity     float64           `json:"sentTareQuantity"`
	SentNetQuantity      float64           `json:"sentNetQuantity"`
	SentTime             *time.Time        `json:"sentTime"`
	RecivedGrossQuantity float64           `json:"recivedGrossQuantity"`
	RecivedTareQuantity  float64           `json:"recivedTareQuantity"`
	RecivedNetQuantity   float64           `json:"recivedNetQuantity"`
	RecivedTime          *time.Time        `json:"recivedTime"`
	Shrinkage            *float64          `json:"shrinkage"`
	Status               string            `json:"status"`
	Number               string            `json:"number"`
	CreateDt             time.Time         `json:"createDt"`
	UpdateBy             string            `json:"updateBy"`
	UpdateDt             time.Time         `json:"updateDt"`
	DeleteDt             gorm.DeletedAt    `json:"deleteDt"`
	CreateName           string            `json:"createName"`
	UpdateName           string            `json:"updateName"`

	Product       *ProductView       `json:"product,omitempty"`
	Vehicle       *VehicleView       `json:"vehicle,omitempty"`
	Stockmovement *StockmovementView `json:"stockmovement,omitempty"`
	Retail        *RetailView        `json:"retail,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Purchaseorder *PurchaseorderView `json:"purchaseorder,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
}

func (StockmovementvehicleView) TableName() string {
	return VIEW_STOCKMOVEMENTVEHICLE
}

type InboundView struct {
	ID                   string            `json:"id"`
	WarehouseID          string            `json:"warehouseId"`
	StockmovementID      string            `json:"stockmovementId"`
	ProductID            string            `json:"productId"`
	VehicleID            string            `json:"vehicleId"`
	Type                 StockMovementType `json:"type"`
	Remark               string            `json:"remark"`
	SentGrossQuantity    float64           `json:"sentGrossQuantity"`
	SentTareQuantity     float64           `json:"sentTareQuantity"`
	SentNetQuantity      float64           `json:"sentNetQuantity"`
	SentTime             *time.Time        `json:"sentTime"`
	RecivedGrossQuantity float64           `json:"recivedGrossQuantity"`
	RecivedTareQuantity  float64           `json:"recivedTareQuantity"`
	RecivedNetQuantity   float64           `json:"recivedNetQuantity"`
	RecivedTime          *time.Time        `json:"recivedTime"`
	Shrinkage            *float64          `json:"shrinkage"`
	Status               string            `json:"status"`
	Number               string            `json:"number"`
	CreateDt             time.Time         `json:"createDt"`
	UpdateBy             string            `json:"updateBy"`
	UpdateDt             time.Time         `json:"updateDt"`
	DeleteDt             gorm.DeletedAt    `json:"deleteDt"`
	CreateName           string            `json:"createName"`
	UpdateName           string            `json:"updateName"`

	Warehouse     *WarehouseView     `json:"warehouse,omitempty"`
	Vehicle       *VehicleView       `json:"vehicle,omitempty"`
	Stockmovement *StockmovementView `json:"stockmovement,omitempty"`
	Product       *ProductView       `json:"product,omitempty"`
}

func (InboundView) TableName() string {
	return VIEW_INBOUND
}

type OutboundView struct {
	ID                   string            `json:"id"`
	WarehouseID          string            `json:"warehouseId"`
	StockmovementID      string            `json:"stockmovementId"`
	ProductID            string            `json:"productId"`
	VehicleID            string            `json:"vehicleId"`
	Type                 StockMovementType `json:"type"`
	Remark               string            `json:"remark"`
	SentGrossQuantity    float64           `json:"sentGrossQuantity"`
	SentTareQuantity     float64           `json:"sentTareQuantity"`
	SentNetQuantity      float64           `json:"sentNetQuantity"`
	SentTime             *time.Time        `json:"sentTime"`
	RecivedGrossQuantity float64           `json:"recivedGrossQuantity"`
	RecivedTareQuantity  float64           `json:"recivedTareQuantity"`
	RecivedNetQuantity   float64           `json:"recivedNetQuantity"`
	RecivedTime          *time.Time        `json:"recivedTime"`
	Shrinkage            *float64          `json:"shrinkage"`
	Status               string            `json:"status"`
	Number               string            `json:"number"`
	CreateDt             time.Time         `json:"createDt"`
	UpdateBy             string            `json:"updateBy"`
	UpdateDt             time.Time         `json:"updateDt"`
	DeleteDt             gorm.DeletedAt    `json:"deleteDt"`
	CreateName           string            `json:"createName"`
	UpdateName           string            `json:"updateName"`

	Warehouse     *WarehouseView     `json:"warehouse,omitempty"`
	Vehicle       *VehicleView       `json:"vehicle,omitempty"`
	Stockmovement *StockmovementView `json:"stockmovement,omitempty"`
	Product       *ProductView       `json:"product,omitempty"`
}

func (OutboundView) TableName() string {
	return VIEW_OUTBOUND
}

type StockinView struct {
	ID            string         `json:"id"`
	WarehouseID   string         `json:"warehouseId"`
	ProductID     string         `json:"productId"`
	Remark        string         `json:"remark"`
	GrossQuantity float64        `json:"grossQuantity"`
	TareQuantity  float64        `json:"tareQuantity"`
	NetQuantity   float64        `json:"netQuantity"`
	CreateDt      time.Time      `json:"createDt"`
	UpdateBy      string         `json:"updateBy"`
	UpdateDt      time.Time      `json:"updateDt"`
	DeleteDt      gorm.DeletedAt `json:"deleteDt"`
	CreateName    string         `json:"createName"`
	UpdateName    string         `json:"updateName"`

	Warehouse *WarehouseView `json:"warehouse,omitempty"`
	Product   *ProductView   `json:"product,omitempty"`
}

func (StockinView) TableName() string {
	return VIEW_STOCKIN
}
