package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Ganti dengan Client ID & Client Secret dari Google Cloud Console
var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}
