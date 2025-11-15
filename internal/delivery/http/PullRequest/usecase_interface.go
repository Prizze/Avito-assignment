package pullrequest

import (
	"context"
	"pr-reviewer/internal/domain"
)

type prUC interface {
	CreatePullRequest(ctx context.Context, cr *domain.CreatePullRequest) (*domain.PullRequest, error)
	MergePullRequest(ctx context.Context, prID int) (*domain.PullRequest, error)
}
