package event

import (
	"context"
	"golang.org/x/xerrors"

	"github.com/nhatminhk63j/uetvoting/pkg/auth"
	"github.com/nhatminhk63j/uetvoting/server/middleware/authentication"
)

// Service ...
type Service interface {
	UpsertEvent(ctx context.Context, event *Event) (eventID int, err error)
	GetEventDetail(ctx context.Context, eventID int) (*Event, error)
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
		eventInfo, err := s.repo.GetEventInfo(ctx, event.ID)
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

// GetEventDetail ...
func (s service) GetEventDetail(ctx context.Context, eventID int) (*Event, error) {
	userInfo := authentication.GetUserInfoFromContext(ctx)
	event, err := s.repo.GetEventDetail(ctx, eventID)
	if err != nil {
		return nil, xerrors.Errorf("error getting event detail: %w", err)
	}
	if userInfo.Id != event.CreatedBy {
		return nil, &auth.NoPermissionError{}
	}
	return event, nil
}
