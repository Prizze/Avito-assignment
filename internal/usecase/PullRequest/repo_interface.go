package pullrequest

import (
	"context"
	"pr-reviewer/internal/domain"
)

type PullRequestRepo interface {
	ExistsById(ctx context.Context, id int) (bool, error)
	GetActiveTeamMembersExceptAuthor(ctx context.Context, authorId int) ([]domain.User, error)
	Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	GetById(ctx context.Context, id int) (*domain.PullRequest, error)
	UpdateStatus(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	UpdateAssignedReviewers(ctx context.Context, prID int, oldReviewerID int, newReviewerID int) error
}
