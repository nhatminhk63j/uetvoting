package user

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpsertUser(ctx context.Context, user *User) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

// GetUserByEmail ...
func (r repo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.Where("Email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ErrNotFound{Email: email}
		}
		return nil, xerrors.Errorf("error getting user %s: %w", email, err)
	}
	return &user, nil
}

// UpsertUser ...
func (r repo) UpsertUser(ctx context.Context, user *User) error {
	err := r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user).Error

	if err != nil {
		return xerrors.Errorf("error create or update user: %s: %w", user.Email, err)
	}
	return nil
}
