package auth

import (
	"github.com/jihanlugas/warehouse/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func init() {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     config.GoogleOauth.ClientID,
		ClientSecret: config.GoogleOauth.ClientSecret,
		RedirectURL:  "http://localhost:1323/auth/google/callback", // ex: http://localhost:8080/auth/google/callback
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
