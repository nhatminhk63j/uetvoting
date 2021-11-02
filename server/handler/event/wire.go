//go:build wireinject

package event

import (
	"github.com/google/wire"
	"github.com/nhatminhk63j/uetvoting/pkg/event"
)

func InitializeHandler() *ServiceServer {
	panic(wire.Build(
		NewServiceServer,
		event.InitializeService,
	))
	return nil
}
