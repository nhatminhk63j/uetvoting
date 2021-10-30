package mysql

import (
	"strconv"

	migrateV4 "github.com/golang-migrate/migrate/v4"

	"github.com/nhatminhk63j/uetvoting/pkg/logger"
)

// ForceMigrate ...
func ForceMigrate(databaseURL string, migrationFolder string, strVersion string) {
	m, err := migrateV4.New(migrationFolder, databaseURL)

	if err != nil {
		logger.Fatalf("error create migration: %v", err)
	}

	version, err := strconv.Atoi(strVersion)
	if err != nil {
		logger.Fatalf("error when force migrate: %v", err)
	}

	logger.Infof("force to version: %d", version)

	if err := m.Force(version); err != nil {
		logger.Fatalf("error when force db to version %d: %v", version, err)
	}
}
