package user

import (
	"context"
	"pr-reviewer/internal/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_user_usecase.go -package=mocks

type userUC interface {
	SetUserIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error)
	GetUserPullRequests(ctx context.Context, userID int) ([]domain.PullRequest, error)
}
