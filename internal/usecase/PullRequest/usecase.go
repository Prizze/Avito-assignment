package pullrequest

import (
	"context"
	"fmt"
	"math/rand"
	"pr-reviewer/internal/domain"
	user "pr-reviewer/internal/usecase/User"
	"time"
)

type PullRequestUsecase struct {
	repo     PullRequestRepo
	userRepo user.UserRepo
}

func NewPullRequestUsecase(repo PullRequestRepo, userRepo user.UserRepo) *PullRequestUsecase {
	return &PullRequestUsecase{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (uc *PullRequestUsecase) CreatePullRequest(ctx context.Context, cr *domain.CreatePullRequest) (*domain.PullRequest, error) {
	if err := uc.checkCreatePRConditions(ctx, cr); err != nil {
		return nil, err
	}

	pr := &domain.PullRequest{
		ID:                cr.PullRequestId,
		Name:              cr.Name,
		AuthorID:          cr.AuthorId,
		CreatedAt:         time.Now(),
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []int{},
	}

	teamMembers, err := uc.repo.GetActiveTeamMembersExceptAuthor(ctx, cr.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	n := domain.MaxReviewersNumber
	if len(teamMembers) < 2 {
		n = len(teamMembers)
	}

	rand.Shuffle(len(teamMembers), func(i, j int) { teamMembers[i], teamMembers[j] = teamMembers[j], teamMembers[i] })

	for i := 0; i < n; i++ {
		pr.AssignedReviewers = append(pr.AssignedReviewers, teamMembers[i].ID)
	}

	createdPR, err := uc.repo.Create(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return createdPR, err
}

func (uc *PullRequestUsecase) MergePullRequest(ctx context.Context, prID int) (*domain.PullRequest, error) {
	prExists, err := uc.checkPRIDExists(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to check PR existance: %w", err)
	}
	if !prExists {
		return nil, domain.ErrPullRequestNotFound
	}

	pr, err := uc.repo.GetById(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR by id: %w", err)
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	updatedPR, err := uc.repo.UpdateStatus(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to update PR status: %w", err)
	}

	return updatedPR, nil
}

func (uc *PullRequestUsecase) checkCreatePRConditions(ctx context.Context, cr *domain.CreatePullRequest) error {
	authorExists, err := uc.userRepo.ExistsById(ctx, cr.AuthorId)
	if err != nil {
		return fmt.Errorf("failed to check author existence: %w", err)
	}
	if !authorExists {
		return domain.ErrUserNotFound
	}

	prExists, err := uc.checkPRIDExists(ctx, cr.PullRequestId)
	if err != nil {
		return fmt.Errorf("failed to check PR existance: %w", err)
	}
	if prExists {
		return domain.ErrPullRequestExists
	}
	return nil
}

func (uc *PullRequestUsecase) checkPRIDExists(ctx context.Context, id int) (bool, error) {
	exists, err := uc.repo.ExistsById(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check pull_request existance: %w", err)
	}
	return exists, nil
}
