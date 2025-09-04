package model

import (
	"time"

	"gorm.io/gorm"
)

type PhotoView struct {
	ID          string         `json:"id"`
	ClientName  string         `json:"clientName"`
	ServerName  string         `json:"serverName"`
	RefTable    string         `json:"refTable"`
	Ext         string         `json:"ext"`
	PhotoPath   string         `json:"photoPath"`
	PhotoSize   int64          `json:"photoSize"`
	PhotoWidth  int64          `json:"photoWidth"`
	PhotoHeight int64          `json:"photoHeight"`
	CreateBy    string         `json:"createBy"`
	CreateDt    time.Time      `json:"createDt"`
	UpdateBy    string         `json:"updateBy"`
	UpdateDt    time.Time      `json:"updateDt"`
	DeleteDt    gorm.DeletedAt `json:"deleteDt"`
	CreateName  string         `json:"createName"`
	UpdateName  string         `json:"updateName"`
}

func (PhotoView) TableName() string {
	return VIEW_PHOTO
}

type PhotoincView struct {
	ID         string         `json:"id"`
	RefTable   string         `json:"eefTable"`
	FolderInc  int64          `json:"folderInc"`
	Folder     string         `json:"folder"`
	Running    int64          `json:"running"`
	CreateBy   string         `json:"createBy"`
	CreateDt   time.Time      `json:"createDt"`
	UpdateBy   string         `json:"updateBy"`
	UpdateDt   time.Time      `json:"updateDt"`
	DeleteDt   gorm.DeletedAt `json:"deleteDt"`
	CreateName string         `json:"createName"`
	UpdateName string         `json:"updateName"`
}

func (PhotoincView) TableName() string {
	return VIEW_PHOTOINC
}

type AuditlogView struct {
	ID                     string         `json:"id"`
	LocationID             string         `json:"locationId"`
	WarehouseID            string         `json:"warehouseId"`
	StockmovementvehicleID string         `json:"stockmovementvehicleId"`
	AuditlogType           AuditlogType   `json:"auditlogType"`
	Title                  string         `json:"title"`
	Description            string         `json:"description"`
	Request                string         `json:"request"`
	Response               string         `json:"response"`
	CreateBy               string         `json:"createBy"`
	CreateDt               time.Time      `json:"createDt"`
	UpdateBy               string         `json:"updateBy"`
	UpdateDt               time.Time      `json:"updateDt"`
	DeleteDt               gorm.DeletedAt `json:"deleteDt"`
	CreateName             string         `json:"createName"`
	UpdateName             string         `json:"updateName"`

	Location             *LocationView             `json:"location,omitempty"`
	Warehouse            *WarehouseView            `json:"warehouse,omitempty"`
	Stockmovementvehicle *StockmovementvehicleView `json:"stockmovementvehicle,omitempty"`
}

func (AuditlogView) TableName() string {
	return VIEW_AUDITLOG
}

type UserView struct {
	ID                string         `json:"id"`
	LocationID        string         `json:"locationId"`
	WarehouseID       string         `json:"warehouseId"`
	UserRole          UserRole       `json:"userRole"`
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
	Userproviders []UserproviderView `json:"userproviders,omitempty" gorm:"foreignKey:UserID"`
	Warehouse     *WarehouseView     `json:"warehouse,omitempty"`
	Location      *LocationView      `json:"location,omitempty"`
}

func (UserView) TableName() string {
	return VIEW_USER
}

type UserproviderView struct {
	ID             string         `json:"id"`
	UserID         string         `json:"userId"`
	ProviderName   string         `json:"providerName"`
	ProviderUserID string         `json:"providerUserId"`
	Email          string         `json:"email"`
	Fullname       string         `json:"fullname"`
	CreateBy       string         `json:"createBy"`
	CreateDt       time.Time      `json:"createDt"`
	UpdateBy       string         `json:"updateBy"`
	UpdateDt       time.Time      `json:"updateDt"`
	DeleteDt       gorm.DeletedAt `json:"deleteDt"`
	CreateName     string         `json:"createName"`
	UpdateName     string         `json:"updateName"`

	User *UserView `json:"user,omitempty"`
}

func (UserproviderView) TableName() string {
	return VIEW_USERPROVIDER
}

type UserprivilegeView struct {
	ID            string         `json:"id"`
	UserID        string         `json:"userId"`
	StockIn       bool           `json:"stockIn"`
	TransferOut   bool           `json:"transferOut"`
	TransferIn    bool           `json:"transferIn"`
	Purchaseorder bool           `json:"purchaseorder"`
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
	RetailStatus RetailStatus   `json:"retailStatus"`
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
	ID                  string              `json:"id"`
	CustomerID          string              `json:"customerId"`
	TotalPrice          float64             `json:"totalPrice"`
	TotalPayment        float64             `json:"totalPayment"`
	Outstanding         float64             `json:"outstanding"`
	Number              string              `json:"number"`
	Notes               string              `json:"notes"`
	PurchaseorderStatus PurchaseorderStatus `json:"purchaseorderStatus"`
	CreateBy            string              `json:"createBy"`
	CreateDt            time.Time           `json:"createDt"`
	UpdateBy            string              `json:"updateBy"`
	UpdateDt            time.Time           `json:"updateDt"`
	DeleteDt            gorm.DeletedAt      `json:"deleteDt"`
	CreateName          string              `json:"createName"`
	UpdateName          string              `json:"updateName"`

	Customer              *CustomerView              `json:"customer,omitempty"`
	Purchaseorderproducts []PurchaseorderproductView `json:"purchaseorderproducts,omitempty" gorm:"foreignKey:PurchaseorderID"`
	Transactions          []TransactionView          `json:"transactions,omitempty" gorm:"foreignKey:RelatedID"`
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
	ID                 string             `json:"id"`
	RelatedID          string             `json:"relatedId"`
	TransactionRelated TransactionRelated `json:"transactionRelated"`
	TransactionType    TransactionType    `json:"transactionType"`
	CustomerID         string             `json:"customerId"`
	Amount             float64            `json:"amount"`
	Notes              string             `json:"notes"`
	Number             string             `json:"number"`
	CreateBy           string             `json:"createBy"`
	CreateDt           time.Time          `json:"createDt"`
	UpdateBy           string             `json:"updateBy"`
	UpdateDt           time.Time          `json:"updateDt"`
	DeleteDt           gorm.DeletedAt     `json:"deleteDt"`
	CreateName         string             `json:"createName"`
	UpdateName         string             `json:"updateName"`

	Customer      *CustomerView      `json:"customer,omitempty"`
	Retail        *RetailView        `json:"retail,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Purchaseorder *PurchaseorderView `json:"purchaseorder,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
}

func (TransactionView) TableName() string {
	return VIEW_TRANSACTION
}

type VehicleView struct {
	ID          string         `json:"id"`
	WarehouseID string         `json:"warehouseId"`
	PlateNumber string         `json:"plateNumber"`
	Name        string         `json:"name"`
	Notes       string         `json:"notes"`
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
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Notes      string         `json:"notes"`
	CreateBy   string         `json:"createBy"`
	CreateDt   time.Time      `json:"createDt"`
	UpdateBy   string         `json:"updateBy"`
	UpdateDt   time.Time      `json:"updateDt"`
	DeleteDt   gorm.DeletedAt `json:"deleteDt"`
	CreateName string         `json:"createName"`
	UpdateName string         `json:"updateName"`
}

func (ProductView) TableName() string {
	return VIEW_PRODUCT
}

type LocationView struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Notes      string         `json:"notes"`
	CreateBy   string         `json:"createBy"`
	CreateDt   time.Time      `json:"createDt"`
	UpdateBy   string         `json:"updateBy"`
	UpdateDt   time.Time      `json:"updateDt"`
	DeleteDt   gorm.DeletedAt `json:"deleteDt"`
	CreateName string         `json:"createName"`
	UpdateName string         `json:"updateName"`

	Warehouses []WarehouseView `json:"warehouses,omitempty" gorm:"foreignKey:LocationID"`
}

func (LocationView) TableName() string {
	return VIEW_LOCATION
}

type WarehouseView struct {
	ID                      string         `json:"id"`
	LocationID              string         `json:"locationId"`
	Name                    string         `json:"name"`
	Address                 string         `json:"address"`
	Notes                   string         `json:"notes"`
	PhoneNumber             string         `json:"phoneNumber"`
	IsStockin               bool           `json:"isStockin"`
	IsTransferIn            bool           `json:"isTransferIn"`
	IsTransferOut           bool           `json:"isTransferOut"`
	IsRetail                bool           `json:"isRetail"`
	IsPurchaseorder         bool           `json:"isPurchaseorder"`
	PhotoID                 string         `json:"photoId"`
	PhotoUrl                string         `json:"photoUrl"`
	TotalRunningTransferout float64        `json:"totalRunningTransferout"`
	TotalRunningTransferin  float64        `json:"totalRunningTransferin"`
	CreateBy                string         `json:"createBy"`
	CreateDt                time.Time      `json:"createDt"`
	UpdateBy                string         `json:"updateBy"`
	UpdateDt                time.Time      `json:"updateDt"`
	DeleteDt                gorm.DeletedAt `json:"deleteDt"`
	CreateName              string         `json:"createName"`
	UpdateName              string         `json:"updateName"`

	Stocks                []StockView                `json:"stocks,omitempty" gorm:"foreignKey:WarehouseID"`
	Stocklogs             []StocklogView             `json:"stocklogs,omitempty" gorm:"foreignKey:WarehouseID"`
	Location              *LocationView              `json:"location,omitempty"`
	Warehousedestinations []WarehousedestinationView `json:"warehousedestinations,omitempty" gorm:"foreignKey:FromWarehouseID"`
}

func (WarehouseView) TableName() string {
	return VIEW_WAREHOUSE
}

type WarehousedestinationView struct {
	ID              string         `json:"id"`
	FromLocationID  string         `json:"fromLocationId"`
	ToLocationID    string         `json:"toLocationId"`
	FromWarehouseID string         `json:"fromWarehouseId"`
	ToWarehouseID   string         `json:"toWarehouseId"`
	CreateBy        string         `json:"createBy"`
	CreateDt        time.Time      `json:"createDt"`
	UpdateBy        string         `json:"updateBy"`
	UpdateDt        time.Time      `json:"updateDt"`
	DeleteDt        gorm.DeletedAt `json:"deleteDt"`
	CreateName      string         `json:"createName"`
	UpdateName      string         `json:"updateName"`

	FromLocation  *LocationView  `json:"fromLocation,omitempty" gorm:"foreignKey:FromLocationID;references:ID"`
	ToLocation    *LocationView  `json:"toLocation,omitempty" gorm:"foreignKey:ToLocationID;references:ID"`
	FromWarehouse *WarehouseView `json:"fromWarehouse,omitempty" gorm:"foreignKey:FromWarehouseID;references:ID"`
	ToWarehouse   *WarehouseView `json:"toWarehouse,omitempty" gorm:"foreignKey:ToWarehouseID;references:ID"`
}

func (WarehousedestinationView) TableName() string {
	return VIEW_WAREHOUSEDESTINATION
}

type StockView struct {
	ID          string         `json:"id"`
	LocationID  string         `json:"locationId"`
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

	Location  *LocationView  `json:"location,omitempty"`
	Warehouse *WarehouseView `json:"warehouse,omitempty"`
	Product   *ProductView   `json:"product,omitempty"`
}

func (StockView) TableName() string {
	return VIEW_STOCK
}

type StocklogView struct {
	ID                     string         `json:"id"`
	WarehouseID            string         `json:"warehouseId"`
	StockID                string         `json:"stockId"`
	StockmovementvehicleID string         `json:"stockmovementvehicleId"`
	ProductID              string         `json:"productId"`
	VehicleID              string         `json:"vehicleId"`
	StocklogType           StocklogType   `json:"stocklogType"`
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
	Stockmovementvehicle *StockmovementvehicleView `json:"stockmovementvehicle,omitempty"`
	Product              *ProductView              `json:"product,omitempty"`
	Vehicle              *VehicleView              `json:"vehicle,omitempty"`
}

func (StocklogView) TableName() string {
	return VIEW_STOCKLOG
}

type StockmovementvehicleView struct {
	ID                         string                     `json:"id"`
	FromLocationID             string                     `json:"fromLocationId"`
	ToLocationID               string                     `json:"toLocationId"`
	FromWarehouseID            string                     `json:"fromWarehouseId"`
	ToWarehouseID              string                     `json:"toWarehouseId"`
	ProductID                  string                     `json:"productId"`
	VehicleID                  string                     `json:"vehicleId"`
	RelatedID                  string                     `json:"relatedId"`
	StockmovementvehicleType   StockmovementvehicleType   `json:"stockmovementvehicleType"`
	Notes                      string                     `json:"notes"`
	SentGrossQuantity          float64                    `json:"sentGrossQuantity"`
	SentTareQuantity           float64                    `json:"sentTareQuantity"`
	SentNetQuantity            float64                    `json:"sentNetQuantity"`
	SentTime                   *time.Time                 `json:"sentTime"`
	ReceivedGrossQuantity      float64                    `json:"receivedGrossQuantity"`
	ReceivedTareQuantity       float64                    `json:"receivedTareQuantity"`
	ReceivedNetQuantity        float64                    `json:"receivedNetQuantity"`
	ReceivedTime               *time.Time                 `json:"receivedTime"`
	Shrinkage                  float64                    `json:"shrinkage"`
	UnitPrice                  float64                    `json:"unitPrice"`
	Number                     string                     `json:"number"`
	StockmovementvehicleStatus StockmovementvehicleStatus `json:"stockmovementvehicleStatus"`
	CreateDt                   time.Time                  `json:"createDt"`
	UpdateBy                   string                     `json:"updateBy"`
	UpdateDt                   time.Time                  `json:"updateDt"`
	DeleteDt                   gorm.DeletedAt             `json:"deleteDt"`
	CreateName                 string                     `json:"createName"`
	UpdateName                 string                     `json:"updateName"`

	FromLocation               *LocationView                   `json:"fromLocation,omitempty" gorm:"foreignKey:FromLocationID;references:ID"`
	ToLocation                 *LocationView                   `json:"toLocation,omitempty" gorm:"foreignKey:ToLocationID;references:ID"`
	FromWarehouse              *WarehouseView                  `json:"fromWarehouse,omitempty" gorm:"foreignKey:FromWarehouseID;references:ID"`
	ToWarehouse                *WarehouseView                  `json:"toWarehouse,omitempty" gorm:"foreignKey:ToWarehouseID;references:ID"`
	Product                    *ProductView                    `json:"product,omitempty"`
	Vehicle                    *VehicleView                    `json:"vehicle,omitempty"`
	Retail                     *RetailView                     `json:"retail,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Purchaseorder              *PurchaseorderView              `json:"purchaseorder,omitempty" gorm:"foreignKey:RelatedID;references:ID"`
	Stockmovementvehiclephotos []StockmovementvehiclephotoView `json:"stockmovementvehiclephotos,omitempty" gorm:"foreignKey:StockmovementvehicleID"`
	Auditlogs                  []AuditlogView                  `json:"auditlogs,omitempty" gorm:"foreignKey:StockmovementvehicleID"`
}

func (StockmovementvehicleView) TableName() string {
	return VIEW_STOCKMOVEMENTVEHICLE
}

type StockmovementvehiclephotoView struct {
	ID                     string         `json:"id"`
	WarehouseID            string         `json:"sarehouseId"`
	StockmovementvehicleID string         `json:"stockmovementvehicleId"`
	PhotoID                string         `json:"photoId"`
	PhotoUrl               string         `json:"photoUrl"`
	CreateDt               time.Time      `json:"createDt"`
	UpdateBy               string         `json:"updateBy"`
	UpdateDt               time.Time      `json:"updateDt"`
	DeleteDt               gorm.DeletedAt `json:"deleteDt"`
	CreateName             string         `json:"createName"`
	UpdateName             string         `json:"updateName"`

	Warehouse            *WarehouseView            `json:"warehouse,omitempty"`
	Stockmovementvehicle *StockmovementvehicleView `json:"stockmovementvehicle,omitempty"`
}

func (StockmovementvehiclephotoView) TableName() string {
	return VIEW_STOCKMOVEMENTVEHICLEPHOTO
}
