package user

import (
	"context"
	"pr-reviewer/internal/domain"
)

type userUC interface {
	SetUserIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error)
}
