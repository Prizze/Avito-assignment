package team

import (
	"context"
	"pr-reviewer/internal/domain"
)

type teamUC interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeamByName(ctx context.Context, name string) (*domain.Team, error)
}
