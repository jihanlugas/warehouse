package response

import "github.com/jihanlugas/warehouse/model"

type Init struct {
	User      model.UserView       `json:"user,omitempty"`
	Warehouse *model.WarehouseView `json:"warehouse,omitempty"`
}
