//go:build wireinject

package auth

import (
	"github.com/google/wire"
	"github.com/nhatminhk63j/uetvoting/pkg/jwt"
	"github.com/nhatminhk63j/uetvoting/pkg/user"
)

func InitializeService() Service {
	panic(wire.Build(
		NewService,
		jwt.NewJWTResolver,
		user.InitializeService,
		wire.Bind(new(UserService), new(user.Service)),
	))
	return nil
}
