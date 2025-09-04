package response

import "github.com/jihanlugas/warehouse/model"

type Init struct {
	User      model.UserView       `json:"user,omitempty"`
	Warehouse *model.WarehouseView `json:"warehouse,omitempty"`
}

type GoogleCallback struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}
