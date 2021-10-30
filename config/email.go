package config

// EmailSender ...
type EmailSender struct {
	Address string `name:"email-address" help:"Email address to send notifications" env:"EMAIL_ADDRESS" default:"nhatns.uet@gmail.com"`
	Name    string `name:"email-name" help:"Email name to send notifications" env:"EMAIL_NAME" default:"Uet Voting"`
}
