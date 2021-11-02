package event

import (
	"context"

	"github.com/nhatminhk63j/uetvoting/pkg/auth"
	"github.com/nhatminhk63j/uetvoting/server/middleware/authentication"
)

// Service ...
type Service interface {
	UpsertEvent(ctx context.Context, event *Event) (eventID int, err error)
}

type service struct {
	repo Repository
}

// NewService ...
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// UpsertEvent ...
func (s service) UpsertEvent(ctx context.Context, event *Event) (eventID int, err error) {
	userInfo := authentication.GetUserInfoFromContext(ctx)
	if event.ID > 0 {
		eventInfo, err := s.repo.GetEventByID(ctx, event.ID)
		if err != nil {
			return 0, err
		}
		if userInfo.Id != eventInfo.CreatedBy {
			return 0, &auth.NoPermissionError{}
		}
		event.UpdatedBy = userInfo.Id
	} else {
		event.CreatedBy = userInfo.Id
		event.UpdatedBy = userInfo.Id
	}
	return s.repo.UpsertEvent(ctx, event)
}
