package request

type PageLocation struct {
	Paging
	Name       string `json:"name" form:"name" query:"name"`
	Address    string `json:"address" form:"address" query:"address" validate:""`
	Notes      string `json:"notes" form:"notes" query:"notes" validate:""`
	CreateName string `json:"createName" form:"createName" query:"createName"`
	Preloads   string `json:"preloads" form:"preloads" query:"preloads"`
}
