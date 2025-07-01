package request

type CreateStockmovement struct {
	FromWarehouseID string `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId" validate:"required"`
	ToWarehouseID   string `json:"toWarehouseId" form:"toWarehouseId" query:"toWarehouseId" validate:"required"`
	ProductID       string `json:"productId" form:"productId" query:"productId" validate:"required"`
	Type            string `json:"type" form:"type" query:"type" validate:"required"`
	Remark          string `json:"remark" form:"remark" query:"remark" validate:""`
}

type UpdateStockmovement struct {
	FromWarehouseID string `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId" validate:"required"`
	ToWarehouseID   string `json:"toWarehouseId" form:"toWarehouseId" query:"toWarehouseId" validate:"required"`
	ProductID       string `json:"productId" form:"productId" query:"productId" validate:"required"`
	Type            string `json:"type" form:"type" query:"type" validate:"required"`
	Remark          string `json:"remark" form:"remark" query:"remark" validate:""`
}

type PageStockmovement struct {
	Paging
	FromWarehouseID string `json:"fromWarehouseId" form:"fromWarehouseId" query:"fromWarehouseId"`
	ToWarehouseID   string `json:"toWarehouseId" form:"toWarehouseId" query:"toWarehouseId"`
	ProductID       string `json:"productId" form:"productId" query:"productId"`
	Type            string `json:"type" form:"type" query:"type"`
	Remark          string `json:"remark" form:"remark" query:"remark"`
	CreateName      string `json:"createName" form:"createName" query:"createName"`
	Preloads        string `json:"preloads" form:"preloads" query:"preloads"`
}
