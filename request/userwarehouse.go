package request

type UpdateUserprivilege struct {
	StockIn       bool `json:"stockIn" form:"stockIn" query:"stockIn"`
	TransferOut   bool `json:"transferOut" form:"transferOut" query:"transferOut"`
	TransferIn    bool `json:"transferIn" form:"transferIn" query:"transferIn"`
	PurchaseOrder bool `json:"purchaseOrder" form:"purchaseOrder" query:"purchaseOrder"`
	Retail        bool `json:"retail" form:"retail" query:"retail"`
}
type CreateUserwarehouse struct {
	WarehouseID        string `json:"warehouseId" validate:"required"`
	UserID             string `json:"userId" validate:"required"`
	IsDefaultWarehouse bool   `json:"isDefaultWarehouse" validate:""`
	IsCreator          bool   `json:"isCreator" validate:""`
}

type UpdateUserwarehouse struct {
	WarehouseID        string `json:"warehouseId" validate:"required"`
	UserID             string `json:"userId" validate:"required"`
	IsDefaultWarehouse bool   `json:"isDefaultWarehouse" validate:""`
	IsCreator          bool   `json:"isCreator" validate:""`
}

type PageUserwarehouse struct {
	Paging
	WarehouseID string `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	UserID      string `json:"userId" form:"userId" query:"userId"`
	Preloads    string `json:"preloads" form:"preloads" query:"preloads"`
}
