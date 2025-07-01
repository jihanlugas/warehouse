package request

import "time"

type CreateCustomer struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber" validate:"required"`
	Email       string `json:"email" form:"email" query:"email" validate:""`
	Address     string `json:"address" form:"address" query:"address" validate:""`
}

type UpdateCustomer struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber" validate:"required"`
	Email       string `json:"email" form:"email" query:"email" validate:""`
	Address     string `json:"address" form:"address" query:"address" validate:""`
}

type PageCustomer struct {
	Paging
	Name          string     `json:"name" form:"name" query:"name"`
	PhoneNumber   string     `json:"phoneNumber" form:"phoneNumber" query:"phoneNumber"`
	Email         string     `json:"email" form:"email" query:"email"`
	Address       string     `json:"address" form:"address" query:"address"`
	CreateName    string     `json:"createName" form:"createName" query:"createName"`
	StartCreateDt *time.Time `json:"startCreateDt" form:"startCreateDt" query:"startCreateDt"`
	EndCreateDt   *time.Time `json:"endCreateDt" form:"endCreateDt" query:"endCreateDt"`
	Preloads      string     `json:"preloads" form:"preloads" query:"preloads"`
}
