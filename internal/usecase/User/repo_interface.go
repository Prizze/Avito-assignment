package user

import (
	"context"
	"pr-reviewer/internal/domain"
)

//go:generate mockgen -source repo_interface.go -destination=mocks/mock_user_repo.go -package=mocks

type UserRepo interface {
	ExistsById(ctx context.Context, id int) (bool, error)
	UpdateIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error)
	GetUserPullRequests(ctx context.Context, userID int) ([]domain.PullRequest, error)
}
