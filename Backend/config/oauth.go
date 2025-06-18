package config

import (
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth3 "google.golang.org/api/oauth2/v2"
	"os"
)

func NewGoogleOAuth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:       google.Endpoint.AuthURL,
			DeviceAuthURL: google.Endpoint.DeviceAuthURL,
			TokenURL:      google.Endpoint.TokenURL,
			AuthStyle:     oauth2.AuthStyleInHeader,
		},
		RedirectURL: "http://" + os.Getenv("HOST") + os.Getenv("PORT") + "/v1/auth/google/callback",
		Scopes:      []string{oauth3.UserinfoEmailScope, oauth3.UserinfoProfileScope, oauth3.OpenIDScope},
	}
}
