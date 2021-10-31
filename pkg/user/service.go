package user

import "context"

// Service ...
type Service interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpsertUser(ctx context.Context, user *User) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetUserByEmail ...
func (s service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// UpsertUser ...
func (s service) UpsertUser(ctx context.Context, user *User) error {
	return s.repo.UpsertUser(ctx, user)
}
