package request

type CreateWarehouse struct {
	LocationID string `json:"locationId" form:"locationId" query:"locationId" validate:"required"`
	Name       string `json:"name" form:"name" query:"name" validate:"required"`
	Address    string `json:"address" form:"address" query:"address" validate:""`
	Notes      string `json:"notes" form:"notes" query:"notes" validate:""`
}
type UpdateWarehouse struct {
	LocationID string `json:"locationId" form:"locationId" query:"locationId" validate:"required"`
	Name       string `json:"name" form:"name" query:"name" validate:"required"`
	Address    string `json:"address" form:"address" query:"address" validate:""`
	Notes      string `json:"notes" form:"notes" query:"notes" validate:""`
}

type PageWarehouse struct {
	Paging
	LocationID string `json:"locationId" form:"locationId" query:"locationId" validate:""`
	Name       string `json:"name" form:"name" query:"name"`
	Address    string `json:"address" form:"address" query:"address" validate:""`
	Notes      string `json:"notes" form:"notes" query:"notes" validate:""`
	CreateName string `json:"createName" form:"createName" query:"createName"`
	Preloads   string `json:"preloads" form:"preloads" query:"preloads"`
}
