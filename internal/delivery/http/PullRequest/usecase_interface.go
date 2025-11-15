package pullrequest

import (
	"context"
	"pr-reviewer/internal/domain"
)

//go:generate mockgen -source usecase_interface.go -destination=mocks/mock_pullrequest_usecase.go -package=mocks

type prUC interface {
	CreatePullRequest(ctx context.Context, cr *domain.CreatePullRequest) (*domain.PullRequest, error)
	MergePullRequest(ctx context.Context, prID int) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, reas *domain.ReassingReviewer) (*domain.PullRequest, int, error)
}
