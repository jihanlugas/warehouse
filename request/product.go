package request

import "time"

type CreateProduct struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type UpdateProduct struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	Description string `json:"description" form:"description" query:"description" validate:"required"`
}

type PageProduct struct {
	Paging
	Name          string     `json:"name" form:"name" query:"name"`
	Description   string     `json:"description" form:"description" query:"description"`
	CreateName    string     `json:"createName" form:"createName" query:"createName"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
