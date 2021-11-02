//go:build wireinject

package event

import (
	"github.com/google/wire"
	"github.com/nhatminhk63j/uetvoting/pkg/db"
)

func InitializeRepo() Repository {
	panic(wire.Build(
		NewRepository,
		db.InitializeDatabase,
	))
	return nil
}

func InitializeService() Service {
	panic(wire.Build(
		NewService,
		InitializeRepo,
	))
}
