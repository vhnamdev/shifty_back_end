package configs

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConfig struct {
	GoogleLoginConfig oauth2.Config
}

func NewGoogleConfig(clientID, clientSecret, redirectURL string) *GoogleConfig {
	return &GoogleConfig{
		GoogleLoginConfig: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  "postmessage",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}
