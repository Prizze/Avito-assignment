package pullrequest

import (
	"context"
	"fmt"
	"math/rand"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/pkg/logger"
	user "pr-reviewer/internal/usecase/User"
	"slices"
	"time"
)

type PullRequestUsecase struct {
	repo     PullRequestRepo
	userRepo user.UserRepo
	logger   logger.Logger
}

func NewPullRequestUsecase(repo PullRequestRepo, userRepo user.UserRepo, logger logger.Logger) *PullRequestUsecase {
	return &PullRequestUsecase{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (uc *PullRequestUsecase) CreatePullRequest(ctx context.Context, cr *domain.CreatePullRequest) (*domain.PullRequest, error) {
	if err := uc.checkCreatePRConditions(ctx, cr.AuthorId, cr.PullRequestId); err != nil {
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
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "authorID": cr.AuthorId}).Error("PR usecase: failed to get active members")
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
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "prID": pr.ID}).Error("PR usecase: failed to create pull_request")
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
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "prID": prID}).Error("PR usecase: failed to get pull_request by id")
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
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "prID": pr.ID, "status": pr.Status}).Error("PR usecase: failed to update status")
		return nil, fmt.Errorf("failed to update PR status: %w", err)
	}

	return updatedPR, nil
}

func (uc *PullRequestUsecase) ReassignReviewer(ctx context.Context, reas *domain.ReassingReviewer) (*domain.PullRequest, int, error) {
	userExists, err := uc.userRepo.ExistsById(ctx, reas.UserID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "userID": reas.UserID}).Error("PR usecase: failed to check user existence")
		return nil, 0, fmt.Errorf("failed to check author existence: %w", err)
	}
	if !userExists {
		return nil, 0, domain.ErrUserNotFound
	}

	prExists, err := uc.checkPRIDExists(ctx, reas.PullRequestID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check PR existance: %w", err)
	}
	if !prExists {
		return nil, 0, domain.ErrPullRequestNotFound
	}

	pr, err := uc.repo.GetById(ctx, reas.PullRequestID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "prID": reas.PullRequestID}).Error("PR usecase: failed to get pull_request by id")
		return nil, 0, fmt.Errorf("failed to get pull_request: %w", err)
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, 0, domain.ErrPullRequestIsMerged
	}

	idx := slices.IndexFunc(pr.AssignedReviewers, func(id int) bool {
		return id == reas.UserID
	})
	if idx == -1 {
		return nil, 0, domain.ErrNotAssigned
	}

	candidates, err := uc.repo.GetActiveTeamMembersExceptAuthor(ctx, pr.AuthorID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "authorID": pr.AuthorID}).Error("PR usecase: failed to get active members")
		return nil, 0, fmt.Errorf("failed to get team members: %w", err)
	}

	filteredCandidates := make([]domain.User, 0)
	for _, u := range candidates {
		if !slices.Contains(pr.AssignedReviewers, u.ID) {
			filteredCandidates = append(filteredCandidates, u)
		}
	}

	if len(filteredCandidates) == 0 {
		return nil, 0, domain.ErrNoAvailableCandidats
	}

	newReviewer := filteredCandidates[rand.Intn(len(filteredCandidates))]

	pr.AssignedReviewers[idx] = newReviewer.ID

	err = uc.repo.UpdateAssignedReviewers(ctx, pr.ID, reas.UserID, newReviewer.ID)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{
			"err": err.Error(), "prID": pr.ID, "old_reviewer": reas.UserID, "new_reviewer": newReviewer.ID}).
			Error("PR usecase: failed to update assigned reviewers")
		return nil, 0, fmt.Errorf("failed to update assigned reviewers: %w", err)
	}

	return pr, newReviewer.ID, nil
}

func (uc *PullRequestUsecase) checkCreatePRConditions(ctx context.Context, uid int, prid int) error {
	authorExists, err := uc.userRepo.ExistsById(ctx, uid)
	if err != nil {
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "userID": uid}).Error("PR usecase: failed to check user existence")
		return fmt.Errorf("failed to check author existence: %w", err)
	}
	if !authorExists {
		return domain.ErrUserNotFound
	}

	prExists, err := uc.checkPRIDExists(ctx, prid)
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
		uc.logger.WithFields(logger.LoggerFields{"err": err.Error(), "prID": id}).Error("PR usecase: failed to check pr existance")
		return false, fmt.Errorf("failed to check pull_request existance: %w", err)
	}
	return exists, nil
}
