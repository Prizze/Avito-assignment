package user

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"
)

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) SetUserIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error) {
	exists, err := uc.checkUserIDExists(ctx, set.ID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, domain.ErrUserNotFound
	}

	updatedUser, err := uc.repo.UpdateIsActive(ctx, set)
	if err != nil {
		return nil, fmt.Errorf("failed to update_is_activeL %w", err)
	}
	return updatedUser, nil

}

func (uc *UserUsecase) GetUserPullRequests(ctx context.Context, userID int) ([]domain.PullRequest, error) {
	exists, err := uc.checkUserIDExists(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, domain.ErrUserNotFound
	}

	userPRs, err := uc.repo.GetUserPullRequests(ctx, userID)
	if err != nil {
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
