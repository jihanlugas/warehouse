package request

import "time"

type CreateTransaction struct {
	RelatedID   string  `json:"relatedId" form:"relatedId" query:"relatedId" validate:"required"`
	RelatedType string  `json:"relatedType" form:"relatedType" query:"relatedType" validate:"required"`
	Type        string  `json:"type" form:"type" query:"type" validate:"required"`
	Amount      float64 `json:"amount" form:"amount" query:"amount" validate:"required"`
	Notes       string  `json:"notes" form:"notes" query:"notes" validate:""`
}

type UpdateTransaction struct {
	Notes string `json:"notes" form:"notes" query:"notes" validate:""`
}

type PageTransaction struct {
	Paging
	CustomerID    string     `json:"customerId" form:"customerId" query:"customerId"`
	RelatedID     string     `json:"relatedId" form:"relatedId" query:"relatedId"`
	RelatedType   string     `json:"relatedType" form:"relatedType" query:"relatedType"`
	CreateName    string     `json:"createName" form:"createName" query:"createName"`
	StartAmount   *float64   `json:"startAmount" form:"startAmount" query:"startAmount"`
	EndAmount     *float64   `json:"endAmount" form:"endAmount" query:"endAmount"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Notes         string     `json:"notes" form:"notes" query:"notes"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
