package mysql

import (
	migrateV4 "github.com/golang-migrate/migrate/v4"

	"github.com/nhatminhk63j/uetvoting/pkg/logger"
)

// MigrateUp migrate db to latest version.
func MigrateUp(databaseURL string, migrationFolder string) {
	m, err := migrateV4.New(migrationFolder, databaseURL)
	if err != nil {
		logger.Fatalf("error create migration: %v", err)
	}

	logger.Info("migrate up completed")

	if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
		logger.Fatalf("error when migrate up: %v", err)
	}
}
