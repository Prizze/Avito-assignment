package team

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"
)

type TeamUsecase struct {
	repo teamRepo
}

func NewTeamUsecase(repo teamRepo) *TeamUsecase {
	return &TeamUsecase{
		repo: repo,
	}
}

func (uc *TeamUsecase) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	exists, err := uc.checkTeamNameExists(ctx, team.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrTeamExists
	}

	createdTeam, err := uc.repo.Create(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return createdTeam, nil
}

func (uc *TeamUsecase) checkTeamNameExists(ctx context.Context, name string) (bool, error) {
	exists, err := uc.repo.ExistsByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check team_name existance: %w", err)
	}
	return exists, err
}

func (uc *TeamUsecase) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	exists, err := uc.checkTeamNameExists(ctx, name)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, domain.ErrTeamNotFound
	}

	team, err := uc.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	return team, nil
}
