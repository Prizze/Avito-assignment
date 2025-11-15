package team

import (
	"context"
	"pr-reviewer/internal/domain"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_team_repo.go -package=mocks

type teamRepo interface {
	ExistsByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetByName(ctx context.Context, name string) (*domain.Team, error)
}
