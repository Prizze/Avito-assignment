package team

import (
	"context"
	"pr-reviewer/internal/domain"
)

type teamRepo interface {
	ExistsByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetByName(ctx context.Context, name string) (*domain.Team, error)
}
