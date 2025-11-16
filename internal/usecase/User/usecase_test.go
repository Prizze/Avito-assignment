package user

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"
	mocksLogger "pr-reviewer/internal/pkg/logger/mocks"
	mockRepo "pr-reviewer/internal/usecase/User/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_SetUserIsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRepo.NewMockUserRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)

	uc := &UserUsecase{repo: repo, logger: logger}

	ctx := context.Background()
	set := &domain.SetUserIsActive{ID: 1, IsActive: true}

	t.Run("checkUserIDExists error", func(t *testing.T) {
		repo.EXPECT().ExistsById(ctx, set.ID).Return(false, fmt.Errorf("db error"))

		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("User usecase: check user_id failed")

		user, err := uc.SetUserIsActive(ctx, set)
		assert.Nil(t, user)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().ExistsById(ctx, set.ID).Return(false, nil)

		user, err := uc.SetUserIsActive(ctx, set)
		assert.Nil(t, user)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("update is_active error", func(t *testing.T) {
		repo.EXPECT().ExistsById(ctx, set.ID).Return(true, nil)
		repo.EXPECT().UpdateIsActive(ctx, set).Return(nil, fmt.Errorf("update failed"))

		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("User usecase: update is_active failed")

		user, err := uc.SetUserIsActive(ctx, set)
		assert.Nil(t, user)
		assert.ErrorContains(t, err, "update failed")
	})

	t.Run("ok", func(t *testing.T) {
		updated := &domain.User{ID: set.ID, IsActive: set.IsActive}
		repo.EXPECT().ExistsById(ctx, set.ID).Return(true, nil)
		repo.EXPECT().UpdateIsActive(ctx, set).Return(updated, nil)

		user, err := uc.SetUserIsActive(ctx, set)
		assert.NoError(t, err)
		assert.Equal(t, updated, user)
	})
}

func TestUserUsecase_GetUserPullRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockRepo.NewMockUserRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)

	uc := &UserUsecase{repo: repo, logger: logger}

	ctx := context.Background()
	userID := 123

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().ExistsById(ctx, userID).Return(false, nil)

		prs, err := uc.GetUserPullRequests(ctx, userID)
		assert.Nil(t, prs)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("repo GetUserPullRequests error", func(t *testing.T) {
		repo.EXPECT().ExistsById(ctx, userID).Return(true, nil)
		repo.EXPECT().GetUserPullRequests(ctx, userID).Return(nil, fmt.Errorf("query failed"))
		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("User usecase: get user prs failed")

		prs, err := uc.GetUserPullRequests(ctx, userID)
		assert.Nil(t, prs)
		assert.ErrorContains(t, err, "query failed")
	})

	t.Run("ok", func(t *testing.T) {
		userPRs := []domain.PullRequest{
			{ID: 1, Name: "Fix bug"},
			{ID: 2, Name: "Add feature"},
		}

		repo.EXPECT().ExistsById(ctx, userID).Return(true, nil)
		repo.EXPECT().GetUserPullRequests(ctx, userID).Return(userPRs, nil)

		prs, err := uc.GetUserPullRequests(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, userPRs, prs)
		assert.Len(t, prs, 2)
	})

	t.Run("ok, but pull_requests are empty", func(t *testing.T) {
		userPRs := []domain.PullRequest{}

		repo.EXPECT().ExistsById(ctx, userID).Return(true, nil)
		repo.EXPECT().GetUserPullRequests(ctx, userID).Return(userPRs, nil)

		prs, err := uc.GetUserPullRequests(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, userPRs, prs)
		assert.Len(t, prs, 0)
	})
}
