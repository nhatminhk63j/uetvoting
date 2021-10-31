package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/nhatminhk63j/uetvoting/cmd/mysql"
	"github.com/nhatminhk63j/uetvoting/config"
	"github.com/nhatminhk63j/uetvoting/pkg/logger"
	"github.com/nhatminhk63j/uetvoting/server"
)

func main() {
	var cfg config.AppConfig
	kongCtx := kong.Parse(&cfg)

	_, err := logger.InitLogger(cfg.LogConfig, logger.LoggerBackendZap)
	if err != nil {
		panic(err)
	}

	databaseURL := fmt.Sprintf("mysql://%s", cfg.MysqlCfg.ToURI())
	switch kongCtx.Command() {
	// start server
	case "server":
		server.Serve(&cfg)
	// migrate up
	case "migrate <command>":
		switch cfg.Migrate.Command {
		case "up":
			mysql.MigrateUp(databaseURL, cfg.MigrationFolder)
		}
	case "migrate <command> <option>":
		switch cfg.Migrate.Command {
		case "create":
			mysql.CreateNewMigration(cfg.MigrationFolder, cfg.Migrate.Option)
		case "down":
			mysql.MigrateDown(databaseURL, cfg.MigrationFolder, cfg.Migrate.Option)
		case "force":
			mysql.ForceMigrate(databaseURL, cfg.MigrationFolder, cfg.Migrate.Option)
		}
	default:
		logger.Fatalf("unknown command: %v", kongCtx.Command())
	}
}
