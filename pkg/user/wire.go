//go:build wireinject

package user

import (
	"github.com/google/wire"

	"github.com/nhatminhk63j/uetvoting/pkg/db"
)

func InitializeRepository() Repository {
	panic(wire.Build(
		NewRepository,
		db.InitializeDatabase,
	))
	return nil
}

func InitializeService() Service {
	panic(wire.Build(
		NewService,
		InitializeRepository,
	))
	return nil
}
