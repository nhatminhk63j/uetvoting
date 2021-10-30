package mysql

import (
	"strconv"

	migrateV4 "github.com/golang-migrate/migrate/v4"

	"github.com/nhatminhk63j/uetvoting/pkg/logger"
)

// MigrateDown ...
func MigrateDown(databaseURL string, migrationFolder string, strVersion string) {
	m, err := migrateV4.New(migrationFolder, databaseURL)

	if err != nil {
		logger.Fatalf("error create migration: %v", err)
	}

	version, err := strconv.Atoi(strVersion)
	if err != nil {
		logger.Fatalf("error when migrate down: %v", err)
	}

	logger.Infof("migration down %d", version)
	if err := m.Steps(-version); err != nil {
		logger.Fatalf("error when migration down: %v", err.Error())
	}
}
