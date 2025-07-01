package request

type PageUserprivilege struct {
	Paging
	UserID     string `json:"userId" form:"userId" query:"userId"`
	CreateName string `json:"createName" form:"createName" query:"createName"`
	Preloads   string `json:"preloads" form:"preloads" query:"preloads"`
}
