package user

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/pkg/logger"
)

type UserUsecase struct {
	repo   UserRepo
	logger logger.Logger
}

func NewUserUsecase(repo UserRepo, logger logger.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) SetUserIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error) {
	exists, err := uc.checkUserIDExists(ctx, set.ID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "userID": set.ID, "isActive": set.IsActive}).Error("User usecase: check user_id failed")
		return nil, err
	}

	if !exists {
		return nil, domain.ErrUserNotFound
	}

	updatedUser, err := uc.repo.UpdateIsActive(ctx, set)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "userID": set.ID, "isActive": set.IsActive}).Error("User usecase: update is_active failed")
		return nil, fmt.Errorf("failed to update_is_active %w", err)
	}
	return updatedUser, nil

}

func (uc *UserUsecase) GetUserPullRequests(ctx context.Context, userID int) ([]domain.PullRequest, error) {
	exists, err := uc.checkUserIDExists(ctx, userID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "userID": userID}).Error("User usecase: check user_id failed")
		return nil, err
	}

	if !exists {
		return nil, domain.ErrUserNotFound
	}

	userPRs, err := uc.repo.GetUserPullRequests(ctx, userID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "userID": userID}).Error("User usecase: get user prs failed")
		return nil, fmt.Errorf("failed to get user pull_requests: %w", err)
	}

	return userPRs, nil
}

func (uc *UserUsecase) checkUserIDExists(ctx context.Context, id int) (bool, error) {
	exists, err := uc.repo.ExistsById(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check user existance: %w", err)
	}
	return exists, nil
}
