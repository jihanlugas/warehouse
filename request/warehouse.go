package request

type CreateWarehouse struct {
	Name     string `json:"name" form:"name" query:"name" validate:"required"`
	Location string `json:"location" form:"location" query:"location" validate:""`
}
type UpdateWarehouse struct {
	Name     string `json:"name" form:"name" query:"name" validate:"required"`
	Location string `json:"location" form:"location" query:"location" validate:""`
}

type PageWarehouse struct {
	Paging
	Name       string `json:"name" form:"name" query:"name"`
	Location   string `json:"location" form:"location" query:"location"`
	IsStock    *bool  `json:"isStock" form:"isStock" query:"isStock"`
	CreateName string `json:"createName" form:"createName" query:"createName"`
	Preloads   string `json:"preloads" form:"preloads" query:"preloads"`
}
