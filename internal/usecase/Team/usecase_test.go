package team

import (
	"context"
	"fmt"
	"testing"

	"pr-reviewer/internal/domain"
	mocksLogger "pr-reviewer/internal/pkg/logger/mocks"
	mockRepo "pr-reviewer/internal/usecase/Team/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTeamUsecase_CreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRepo.NewMockteamRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)

	uc := &TeamUsecase{repo: repo, logger: logger}

	ctx := context.Background()
	team := &domain.Team{
		Name: "Team A",
		Members: []domain.TeamMember{
			{UserID: 123, Username: "garry", IsActive: true},
			{UserID: 321, Username: "larry", IsActive: false},
		},
	}

	t.Run("checkTeamNameExists error", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, team.Name).Return(false, fmt.Errorf("db error"))
		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("Team usecase: check team_name existance failed")

		created, err := uc.CreateTeam(ctx, team)
		assert.Nil(t, created)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("team already exists", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, team.Name).Return(true, nil)

		created, err := uc.CreateTeam(ctx, team)
		assert.Nil(t, created)
		assert.Equal(t, domain.ErrTeamExists, err)
	})

	t.Run("repo Create error", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, team.Name).Return(false, nil)
		repo.EXPECT().Create(ctx, team).Return(nil, fmt.Errorf("insert failed"))

		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("Team usecase: create team failed")

		created, err := uc.CreateTeam(ctx, team)
		assert.Nil(t, created)
		assert.ErrorContains(t, err, "insert failed")
	})

	t.Run("ok", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, team.Name).Return(false, nil)
		repo.EXPECT().Create(ctx, team).Return(team, nil)

		created, err := uc.CreateTeam(ctx, team)
		assert.NoError(t, err)
		assert.Equal(t, team, created)
	})
}

func TestTeamUsecase_GetTeamByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRepo.NewMockteamRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)

	uc := &TeamUsecase{repo: repo, logger: logger}

	ctx := context.Background()

	teamName := "Team A"
	team := &domain.Team{
		Name: "Team A",
		Members: []domain.TeamMember{
			{UserID: 123, Username: "garry", IsActive: true},
			{UserID: 321, Username: "larry", IsActive: false},
		},
	}

	t.Run("checkTeamNameExists error", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, teamName).Return(false, fmt.Errorf("db error"))
		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("Team usecase: check team_name existance failed")

		result, err := uc.GetTeamByName(ctx, teamName)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("team not found", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, teamName).Return(false, nil)

		result, err := uc.GetTeamByName(ctx, teamName)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrTeamNotFound, err)
	})

	t.Run("repo GetByName error", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, teamName).Return(true, nil)
		repo.EXPECT().GetByName(ctx, teamName).Return(nil, fmt.Errorf("query failed"))
		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("Team usecase: create team failed")

		result, err := uc.GetTeamByName(ctx, teamName)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "query failed")
	})

	t.Run("ok", func(t *testing.T) {
		repo.EXPECT().ExistsByName(ctx, teamName).Return(true, nil)
		repo.EXPECT().GetByName(ctx, teamName).Return(team, nil)

		result, err := uc.GetTeamByName(ctx, teamName)
		assert.NoError(t, err)
		assert.Equal(t, team, result)
	})
}
