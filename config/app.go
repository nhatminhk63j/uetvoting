package config

import (
	"github.com/nhatminhk63j/uetvoting/pkg/logger"
)

type AppMode string

const (
	// DevelopmentMode indicates that app is in development, so more verbose log.
	DevelopmentMode AppMode = "development"
	// ProductionMode indicates that app is running in production.
	ProductionMode AppMode = "production"
)

type AppConfig struct {
	Migrate struct {
		Command string `kong:"arg,name:'command',enum:'up,down,force,create'"`
		Option  string `kong:"arg,optional,name:'option'"`
	} `kong:"cmd"`

	Server struct{} `kong:"cmd"`

	MigrationFolder string `name:"reviews-folder" help:"Migration folder" env:"MIGRATION_FOLDER" default:"file://migration"`

	AppMode   AppMode              `name:"app-mode" help:"App mode" env:"APP_MODE" default:"production"`
	LogConfig logger.Configuration `kong:"help:'Logger config',embed"`

	MysqlCfg MysqlConfig `kong:"embed,help:'Mysql config'"`

	SendGridKey  string      `name:"sendgrid-key" help:"Send Gird Key" env:"SENDGRID_KEY" default:"teko_key"`
	EmailSender  EmailSender `kong:"help:'Email sender config',embed"`
	ClientDomain string      `name:"client-domain" help:"Send Gird Key" env:"CLIENT_DOMAIN" default:"https://uet-voting.com"`

	GRPCPort  int    `name:"grpc-port" help:"GRPC port" env:"GRPC_PORT" default:"8081"`
	HTTPPort  int    `name:"http-port" help:"HTTP port" env:"HTTP_PORT" default:"8080"`
	SentryDsn string `name:"sentry-dsn" help:"Sentry Dsn" env:"SENTRY_DSN"`
}
