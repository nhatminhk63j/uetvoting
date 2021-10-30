package mysql

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/nhatminhk63j/uetvoting/pkg/logger"
)

const versionTimeFormat = "20060102150405"

// CreateNewMigration with up and down options.
func CreateNewMigration(migrationFolder string, name string) {
	folder := strings.ReplaceAll(migrationFolder, "file://", "")
	now := time.Now()
	ver := now.Format(versionTimeFormat)

	up := fmt.Sprintf("%s/%s_%s.up.sql", folder, ver, name)
	down := fmt.Sprintf("%s/%s_%s.down.sql", folder, ver, name)

	logger.Infof("create migration: %s", name)

	if err := ioutil.WriteFile(up, []byte{}, 0644); err != nil {
		logger.Fatalf("create migration up error: %v", err)
	}
	if err := ioutil.WriteFile(down, []byte{}, 0644); err != nil {
		logger.Fatalf("create migration down error: %v", err)
	}
}
