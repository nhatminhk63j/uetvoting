//go:build wireinject

package db

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	wire.Build(newDatabase)
	return nil
}
