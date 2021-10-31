package config

import "github.com/alecthomas/kong"

// AuthConfig ...
type AuthConfig struct {
	OauthClientID string `name:"oauth-client-id" help:"Oauth Client ID" env:"OAUTH_CLIENT_ID" default:""`
}

// LoadAuthConfig ...
func LoadAuthConfig() *AuthConfig {
	var config AuthConfig
	parser, err := kong.New(&config)
	if err != nil {
		panic(err)
	}
	_, err = parser.Parse([]string{})
	if err != nil {
		panic(err)
	}
	return &config
}
