package user

import (
	"context"
	"pr-reviewer/internal/domain"
)

type UserRepo interface {
	ExistsById(ctx context.Context, id int) (bool, error)
	UpdateIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error)
}
