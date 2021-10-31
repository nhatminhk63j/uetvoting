package config

import "github.com/alecthomas/kong"

// JWTConfig ...
type JWTConfig struct {
	SecretKey    string `name:"secret-key" help:"Secret key" env:"JWT_SECRET_KEY" default:"secret"`
	TimeDuration int    `name:"time-duration" help:"Time Duration" env:"JWT_TIME_DURATION" default:"60"`
}

// LoadJWTConfig ...
func LoadJWTConfig() *JWTConfig {
	var config JWTConfig
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
