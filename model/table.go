package model

import (
	"time"

	"gorm.io/gorm"
)

type PhotoRef string
type UserRole string
type StockmovementvehicleType string
type StocklogType string
type TransactionType string
type TransactionRelated string
type RetailStatus string
type PurchaseorderStatus string
type StockmovementvehicleStatus string

type TransactionStatus string

const (
	PhotoRefStockmovementvehiclephoto PhotoRef = "stockmovementvehiclephoto"
)

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleOperator UserRole = "OPERATOR"
)

const (
	StockmovementvehicleTypeTransfer      StockmovementvehicleType = "TRANSFER"
	StockmovementvehicleTypeIn            StockmovementvehicleType = "IN"
	StockmovementvehicleTypePurchaseorder StockmovementvehicleType = "PURCHASE_ORDER"
	StockmovementvehicleTypeRetail        StockmovementvehicleType = "RETAIL"
)

const (
	StocklogTypeIn  StocklogType = "IN"
	StocklogTypeOut StocklogType = "OUT"
)

const (
	TransactionTypeInvoice TransactionType = "INVOICE"
	TransactionTypePayment TransactionType = "PAYMENT"
)

const (
	TransactionRelatedRetail        TransactionRelated = "RETAIL"
	TransactionRelatedPurchaseorder TransactionRelated = "PURCHASE_ORDER"
)

const (
	RetailStatusOpen  RetailStatus = "OPEN"
	RetailStatusClose RetailStatus = "CLOSE"
)

const (
	PurchaseorderStatusOpen  PurchaseorderStatus = "OPEN"
	PurchaseorderStatusClose PurchaseorderStatus = "CLOSE"
)
const (
	StockmovementvehicleStatusLoading   StockmovementvehicleStatus = "LOADING"
	StockmovementvehicleStatusInTransit StockmovementvehicleStatus = "IN_TRANSIT"
	StockmovementvehicleStatusUnloading StockmovementvehicleStatus = "UNLOADING"
	StockmovementvehicleStatusCompleted StockmovementvehicleStatus = "COMPLETED"
)

type Photo struct {
	ID          string         `gorm:"primaryKey"`
	ClientName  string         `gorm:"not null"`
	ServerName  string         `gorm:"not null"`
	RefTable    string         `gorm:"not null"`
	Ext         string         `gorm:"not null"`
	PhotoPath   string         `gorm:"not null"`
	PhotoSize   int64          `gorm:"not null"`
	PhotoWidth  int64          `gorm:"not null"`
	PhotoHeight int64          `gorm:"not null"`
	CreateBy    string         `gorm:"not null"`
	CreateDt    time.Time      `gorm:"not null"`
	UpdateBy    string         `gorm:"not null"`
	UpdateDt    time.Time      `gorm:"not null"`
	DeleteDt    gorm.DeletedAt `gorm:"null"`
}

type Photoinc struct {
	ID        string         `gorm:"primaryKey"`
	RefTable  string         `gorm:"not null"`
	FolderInc int64          `gorm:"not null"`
	Folder    string         `gorm:"not null"`
	Running   int64          `gorm:"not null"`
	CreateBy  string         `gorm:"not null"`
	CreateDt  time.Time      `gorm:"not null"`
	UpdateBy  string         `gorm:"not null"`
	UpdateDt  time.Time      `gorm:"not null"`
	DeleteDt  gorm.DeletedAt `gorm:"null"`
}

type User struct {
	ID                string         `gorm:"primaryKey"`
	LocationID        string         `gorm:"not null"`
	WarehouseID       string         `gorm:"not null"`
	UserRole          UserRole       `gorm:"not null"`
	Email             string         `gorm:"not null"`
	Username          string         `gorm:"not null"`
	PhoneNumber       string         `gorm:"not null"`
	Address           string         `gorm:"not null"`
	Fullname          string         `gorm:"not null"`
	Passwd            string         `gorm:"not null"`
	PassVersion       int            `gorm:"not null"`
	IsActive          bool           `gorm:"not null"`
	PhotoID           string         `gorm:"not null"`
	LastLoginDt       *time.Time     `gorm:"null"`
	BirthDt           *time.Time     `gorm:"null"`
	BirthPlace        string         `gorm:"not null"`
	AccountVerifiedDt *time.Time     `gorm:"null"`
	CreateBy          string         `gorm:"not null"`
	CreateDt          time.Time      `gorm:"not null"`
	UpdateBy          string         `gorm:"not null"`
	UpdateDt          time.Time      `gorm:"not null"`
	DeleteDt          gorm.DeletedAt `gorm:"null"`

	Userprivilege *Userprivilege `gorm:"not null"`
}

type Userprivilege struct {
	ID            string         `gorm:"primaryKey"`
	UserID        string         `gorm:"not null"`
	StockIn       bool           `gorm:"not null"`
	TransferOut   bool           `gorm:"not null"`
	TransferIn    bool           `gorm:"not null"`
	Purchaseorder bool           `gorm:"not null"`
	Retail        bool           `gorm:"not null"`
	CreateBy      string         `gorm:"not null"`
	CreateDt      time.Time      `gorm:"not null"`
	UpdateBy      string         `gorm:"not null"`
	UpdateDt      time.Time      `gorm:"not null"`
	DeleteDt      gorm.DeletedAt `gorm:"null"`
}

type Customer struct {
	ID          string         `gorm:"primaryKey"`
	Name        string         `gorm:"not null"`
	PhoneNumber string         `gorm:"not null"`
	Email       string         `gorm:"not null"`
	Address     string         `gorm:"not null"`
	CreateBy    string         `gorm:"not null"`
	CreateDt    time.Time      `gorm:"not null"`
	UpdateBy    string         `gorm:"not null"`
	UpdateDt    time.Time      `gorm:"not null"`
	DeleteDt    gorm.DeletedAt `gorm:"null"`
}

type Retail struct {
	ID           string         `gorm:"primaryKey"`
	CustomerID   string         `gorm:"not null"`
	Notes        string         `gorm:"not null"`
	Number       string         `gorm:"not null"`
	RetailStatus RetailStatus   `gorm:"not null"`
	CreateBy     string         `gorm:"not null"`
	CreateDt     time.Time      `gorm:"not null"`
	UpdateBy     string         `gorm:"not null"`
	UpdateDt     time.Time      `gorm:"not null"`
	DeleteDt     gorm.DeletedAt `gorm:"null"`
}

type Retailproduct struct {
	ID        string         `gorm:"primaryKey"`
	RetailID  string         `gorm:"not null"`
	ProductID string         `gorm:"not null"`
	UnitPrice float64        `gorm:"not null"`
	CreateBy  string         `gorm:"not null"`
	CreateDt  time.Time      `gorm:"not null"`
	UpdateBy  string         `gorm:"not null"`
	UpdateDt  time.Time      `gorm:"not null"`
	DeleteDt  gorm.DeletedAt `gorm:"null"`
}

type Purchaseorder struct {
	ID                  string              `gorm:"primaryKey"`
	CustomerID          string              `gorm:"not null"`
	Notes               string              `gorm:"not null"`
	Number              string              `gorm:"not null"`
	PurchaseorderStatus PurchaseorderStatus `gorm:"not null"`
	CreateBy            string              `gorm:"not null"`
	CreateDt            time.Time           `gorm:"not null"`
	UpdateBy            string              `gorm:"not null"`
	UpdateDt            time.Time           `gorm:"not null"`
	DeleteDt            gorm.DeletedAt      `gorm:"null"`
}

type Purchaseorderproduct struct {
	ID              string         `gorm:"primaryKey"`
	PurchaseorderID string         `gorm:"not null"`
	ProductID       string         `gorm:"not null"`
	UnitPrice       float64        `gorm:"not null"`
	CreateBy        string         `gorm:"not null"`
	CreateDt        time.Time      `gorm:"not null"`
	UpdateBy        string         `gorm:"not null"`
	UpdateDt        time.Time      `gorm:"not null"`
	DeleteDt        gorm.DeletedAt `gorm:"null"`
}

type Transaction struct {
	ID                 string             `gorm:"primaryKey"`
	RelatedID          string             `gorm:"not null"`
	TransactionRelated TransactionRelated `gorm:"not null"`
	TransactionType    TransactionType    `gorm:"not null"`
	Amount             float64            `gorm:"not null"`
	Notes              string             `gorm:"not null"`
	Number             string             `gorm:"not null"`
	CreateBy           string             `gorm:"not null"`
	CreateDt           time.Time          `gorm:"not null"`
	UpdateBy           string             `gorm:"not null"`
	UpdateDt           time.Time          `gorm:"not null"`
	DeleteDt           gorm.DeletedAt     `gorm:"null"`
}

type Vehicle struct {
	ID          string         `gorm:"primaryKey"`
	WarehouseID string         `gorm:"not null"`
	PlateNumber string         `gorm:"not null"`
	Name        string         `gorm:"not null"`
	Notes       string         `gorm:"not null"`
	NIK         string         `gorm:"not null"`
	DriverName  string         `gorm:"not null"`
	PhoneNumber string         `gorm:"not null"`
	CreateBy    string         `gorm:"not null"`
	CreateDt    time.Time      `gorm:"not null"`
	UpdateBy    string         `gorm:"not null"`
	UpdateDt    time.Time      `gorm:"not null"`
	DeleteDt    gorm.DeletedAt `gorm:"null"`
}

type Product struct {
	ID       string         `gorm:"primaryKey"`
	Name     string         `gorm:"not null"`
	Notes    string         `gorm:"not null"`
	CreateBy string         `gorm:"not null"`
	CreateDt time.Time      `gorm:"not null"`
	UpdateBy string         `gorm:"not null"`
	UpdateDt time.Time      `gorm:"not null"`
	DeleteDt gorm.DeletedAt `gorm:"null"`
}

type Location struct {
	ID       string         `gorm:"primaryKey"`
	Name     string         `gorm:"primaryKey"`
	Notes    string         `gorm:"not null"`
	CreateBy string         `gorm:"not null"`
	CreateDt time.Time      `gorm:"not null"`
	UpdateBy string         `gorm:"not null"`
	UpdateDt time.Time      `gorm:"not null"`
	DeleteDt gorm.DeletedAt `gorm:"null"`
}

type Warehouse struct {
	ID              string         `gorm:"primaryKey"`
	LocationID      string         `gorm:"not null"`
	Name            string         `gorm:"not null"`
	Address         string         `gorm:"not null"`
	Notes           string         `gorm:"not null"`
	PhoneNumber     string         `gorm:"not null"`
	IsStockin       bool           `gorm:"not null"`
	IsTransferIn    bool           `gorm:"not null"`
	IsTransferOut   bool           `gorm:"not null"`
	IsRetail        bool           `gorm:"not null"`
	IsPurchaseorder bool           `gorm:"not null"`
	CreateBy        string         `gorm:"not null"`
	CreateDt        time.Time      `gorm:"not null"`
	UpdateBy        string         `gorm:"not null"`
	UpdateDt        time.Time      `gorm:"not null"`
	DeleteDt        gorm.DeletedAt `gorm:"null"`
}

type Warehousedestination struct {
	ID              string         `gorm:"primaryKey"`
	FromLocationID  string         `gorm:"not null"`
	ToLocationID    string         `gorm:"not null"`
	FromWarehouseID string         `gorm:"not null"`
	ToWarehouseID   string         `gorm:"not null"`
	CreateBy        string         `gorm:"not null"`
	CreateDt        time.Time      `gorm:"not null"`
	UpdateBy        string         `gorm:"not null"`
	UpdateDt        time.Time      `gorm:"not null"`
	DeleteDt        gorm.DeletedAt `gorm:"null"`
}

type Stock struct {
	ID          string         `gorm:"primaryKey"`
	LocationID  string         `gorm:"not null"`
	WarehouseID string         `gorm:"not null"`
	ProductID   string         `gorm:"not null"`
	Quantity    float64        `gorm:"not null"`
	CreateBy    string         `gorm:"not null"`
	CreateDt    time.Time      `gorm:"not null"`
	UpdateBy    string         `gorm:"not null"`
	UpdateDt    time.Time      `gorm:"not null"`
	DeleteDt    gorm.DeletedAt `gorm:"null"`
}

type Stocklog struct {
	ID                     string         `gorm:"primaryKey"`
	WarehouseID            string         `gorm:"not null"`
	StockID                string         `gorm:"not null"`
	StockmovementvehicleID string         `gorm:"not null"`
	ProductID              string         `gorm:"not null"`
	VehicleID              string         `gorm:"not null"`
	StocklogType           StocklogType   `gorm:"not null"`
	GrossQuantity          float64        `gorm:"not null"`
	TareQuantity           float64        `gorm:"not null"`
	NetQuantity            float64        `gorm:"not null"`
	CreateBy               string         `gorm:"not null"`
	CreateDt               time.Time      `gorm:"not null"`
	UpdateBy               string         `gorm:"not null"`
	UpdateDt               time.Time      `gorm:"not null"`
	DeleteDt               gorm.DeletedAt `gorm:"null"`
}

type Stockmovementvehicle struct {
	ID                         string                     `gorm:"primaryKey"`
	FromLocationID             string                     `gorm:"not null"`
	ToLocationID               string                     `gorm:"not null"`
	FromWarehouseID            string                     `gorm:"not null"`
	ToWarehouseID              string                     `gorm:"not null"`
	ProductID                  string                     `gorm:"not null"`
	VehicleID                  string                     `gorm:"not null"`
	RelatedID                  string                     `gorm:"not null"`
	StockmovementvehicleType   StockmovementvehicleType   `gorm:"not null"`
	Notes                      string                     `gorm:"not null"`
	SentGrossQuantity          float64                    `gorm:"not null"`
	SentTareQuantity           float64                    `gorm:"not null"`
	SentNetQuantity            float64                    `gorm:"not null"`
	SentTime                   *time.Time                 `gorm:"null"`
	ReceivedGrossQuantity      float64                    `gorm:"not null"`
	ReceivedTareQuantity       float64                    `gorm:"not null"`
	ReceivedNetQuantity        float64                    `gorm:"not null"`
	ReceivedTime               *time.Time                 `gorm:"null"`
	Shrinkage                  float64                    `gorm:"not null"`
	Number                     string                     `gorm:"not null"`
	StockmovementvehicleStatus StockmovementvehicleStatus `gorm:"not null"`
	CreateBy                   string                     `gorm:"not null"`
	CreateDt                   time.Time                  `gorm:"not null"`
	UpdateBy                   string                     `gorm:"not null"`
	UpdateDt                   time.Time                  `gorm:"not null"`
	DeleteDt                   gorm.DeletedAt             `gorm:"null"`
}

type Stockmovementvehiclephoto struct {
	ID                     string         `gorm:"primaryKey"`
	WarehouseID            string         `gorm:"not null"`
	StockmovementvehicleID string         `gorm:"not null"`
	PhotoID                string         `json:"photoId"`
	CreateBy               string         `gorm:"not null"`
	CreateDt               time.Time      `gorm:"not null"`
	UpdateBy               string         `gorm:"not null"`
	UpdateDt               time.Time      `gorm:"not null"`
	DeleteDt               gorm.DeletedAt `gorm:"null"`
}
