package request

type CreateStock struct {
	WarehouseID string  `json:"warehouseId" form:"warehouseId" query:"warehouseId" validate:"required"`
	ProductID   string  `json:"productId" form:"productId" query:"productId" validate:"required"`
	Quantity    float64 `json:"quantity" form:"quantity" query:"quantity" validate:"required"`
}

type UpdateStock struct {
	Quantity float64 `json:"quantity" form:"quantity" query:"quantity" validate:""`
}

type PageStock struct {
	Paging
	WarehouseID   string   `json:"warehouseId" form:"warehouseId" query:"warehouseId"`
	ProductID     string   `json:"productId" form:"productId" query:"productId"`
	StartQuantity *float64 `json:"startQuantity" form:"startQuantity" query:"startQuantity"`
	EndQuantity   *float64 `json:"endQuantity" form:"endQuantity" query:"endQuantity"`
	CreateName    string   `json:"createName" form:"createName" query:"createName"`
	Preloads      string   `json:"preloads" form:"preloads" query:"preloads"`
}
