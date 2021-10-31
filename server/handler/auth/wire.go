//go:build wireinject

package auth

import (
	"github.com/google/wire"
	"github.com/nhatminhk63j/uetvoting/pkg/auth"
)

func InitializeHandler() *ServiceServer {
	panic(wire.Build(
		NewServiceServer,
		auth.InitializeService,
	))
	return nil
}
