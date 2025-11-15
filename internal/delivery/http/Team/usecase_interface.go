package team

import (
	"context"
	"pr-reviewer/internal/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_team_usecase.go -package=mocks

type teamUC interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeamByName(ctx context.Context, name string) (*domain.Team, error)
}
