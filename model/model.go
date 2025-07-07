package model

import "time"

const (
	VIEW_USER                 = "users_view"
	VIEW_USERPRIVILEGE        = "userprivileges_view"
	VIEW_CUSTOMER             = "customers_view"
	VIEW_RETAIL               = "retails_view"
	VIEW_RETAILPRODUCT        = "retailproducts_view"
	VIEW_PURCHASEORDER        = "purchaseorders_view"
	VIEW_PURCHASEORDERPRODUCT = "purchaseorderproducts_view"
	VIEW_TRANSACTION          = "transactions_view"
	VIEW_WAREHOUSE            = "wahouses_view"
	VIEW_STOCK                = "stocks_view"
	VIEW_STOCKLOG             = "stocklogs_view"
	VIEW_VEHICLE              = "vehicles_view"
	VIEW_PRODUCT              = "products_view"
	VIEW_STOCKMOVEMENT        = "stockmovements_view"
	VIEW_STOCKMOVEMENTVEHICLE = "stockmovementvehicles_view"

	VIEW_OUTBOUND              = "outbounds_view"
	VIEW_INBOUND               = "inbounds_view"
	VIEW_STOCKIN               = "stockins_view"
	VIEW_DELIVERYRETAIL        = "deliveryretails_view"
	VIEW_DELIVERYPURCHASEORDER = "deliverypurchaseorders_view"
)

type UserLogin struct {
	ExpiredDt       time.Time `json:"expiredDt"`
	UserID          string    `json:"userId"`
	PassVersion     int       `json:"passVersion"`
	WarehouseID     string    `json:"warehouseId"`
	Role            string    `json:"role"`
	UserwarehouseID string    `json:"userwarehouseId"`
}
