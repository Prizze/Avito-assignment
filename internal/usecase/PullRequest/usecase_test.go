package pullrequest

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"
	mocksLogger "pr-reviewer/internal/pkg/logger/mocks"
	mocksRepo "pr-reviewer/internal/usecase/PullRequest/mocks"
	mocksUserRepo "pr-reviewer/internal/usecase/User/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePullRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocksRepo.NewMockPullRequestRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)
	userRepo := mocksUserRepo.NewMockUserRepo(ctrl)

	uc := PullRequestUsecase{repo: repo, logger: logger, userRepo: userRepo}

	ctx := context.Background()
	cr := &domain.CreatePullRequest{
		PullRequestId: 1001,
		Name:          "Test PR",
		AuthorId:      10,
	}

	t.Run("PR created ok", func(t *testing.T) {
		// Участники команды с isActive = true
		members := []domain.User{
			{ID: 11}, {ID: 12}, {ID: 13},
		}

		userRepo.EXPECT().ExistsById(ctx, cr.AuthorId).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, cr.PullRequestId).Return(false, nil)
		repo.EXPECT().GetActiveTeamMembersExceptAuthor(ctx, 10).Return(members, nil)

		repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
				return pr, nil
			},
		)

		pr, err := uc.CreatePullRequest(ctx, cr)

		assert.NoError(t, err)
		assert.NotNil(t, pr)

		// Проверяем что, pr.AssignedReviewers находятся в members
		memberIDs := make([]int, len(members))
		for i, m := range members {
			memberIDs[i] = m.ID
		}

		for _, id := range pr.AssignedReviewers {
			assert.Contains(t, memberIDs, id)
		}

		// Проверяем, что ревьюверов было назначенр =< 2
		assert.LessOrEqual(t, len(pr.AssignedReviewers), domain.MaxReviewersNumber)
	})

	t.Run("PR created, but without reviewers", func(t *testing.T) {
		// Если участников команды с isActive = true, кроме автора, нет
		members := []domain.User{}

		userRepo.EXPECT().ExistsById(ctx, cr.AuthorId).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, cr.PullRequestId).Return(false, nil)
		repo.EXPECT().GetActiveTeamMembersExceptAuthor(ctx, 10).Return(members, nil)

		repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
				return pr, nil
			},
		)

		pr, err := uc.CreatePullRequest(ctx, cr)

		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, len(pr.AssignedReviewers), 0)
	})

}

func TestCheckCreatePRConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocksRepo.NewMockPullRequestRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)
	userRepo := mocksUserRepo.NewMockUserRepo(ctrl)

	uc := PullRequestUsecase{repo: repo, logger: logger, userRepo: userRepo}

	ctx := context.Background()

	t.Run("user exists and PR does not exists - ok", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(gomock.Any(), 10).Return(true, nil)

		repo.EXPECT().ExistsById(gomock.Any(), 5).Return(false, nil)

		err := uc.checkCreatePRConditions(ctx, 10, 5)
		assert.NoError(t, err)
	})

	t.Run("error - user does not exist", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, 11).Return(false, nil)

		err := uc.checkCreatePRConditions(ctx, 11, 5)
		assert.Equal(t, err, domain.ErrUserNotFound)
	})

	t.Run("error - PR already exists", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, 12).Return(true, nil)

		repo.EXPECT().ExistsById(ctx, 7).Return(true, nil)

		err := uc.checkCreatePRConditions(ctx, 12, 7)
		assert.Equal(t, err, domain.ErrPullRequestExists)
	})
}

func TestMergePullRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocksRepo.NewMockPullRequestRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)
	userRepo := mocksUserRepo.NewMockUserRepo(ctrl)

	uc := PullRequestUsecase{repo: repo, logger: logger, userRepo: userRepo}

	prID := 1
	ctx := context.Background()
	now := time.Now()

	t.Run("PR does not exist", func(t *testing.T) {
		repo.EXPECT().ExistsById(gomock.Any(), prID).Return(false, nil)

		pr, err := uc.MergePullRequest(ctx, prID)
		assert.Nil(t, pr)
		assert.Equal(t, domain.ErrPullRequestNotFound, err)
	})

	t.Run("PR already merged", func(t *testing.T) {
		existingPR := &domain.PullRequest{
			ID:     prID,
			Status: domain.PRStatusMerged,
		}

		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(existingPR, nil)

		pr, err := uc.MergePullRequest(ctx, prID)
		assert.NoError(t, err)
		assert.Equal(t, domain.PRStatusMerged, pr.Status)
	})

	t.Run("update status error", func(t *testing.T) {
		prToMerge := &domain.PullRequest{
			ID:     prID,
			Status: domain.PRStatusOpen,
		}

		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(prToMerge, nil)
		repo.EXPECT().UpdateStatus(ctx, gomock.Any()).Return(nil, fmt.Errorf("update failed"))

		logger.EXPECT().WithFields(gomock.Any()).Return(logger)
		logger.EXPECT().Error("PR usecase: failed to update status")

		pr, err := uc.MergePullRequest(ctx, prID)
		assert.Nil(t, pr)
		assert.ErrorContains(t, err, "update failed")
	})

	t.Run("successfully merged", func(t *testing.T) {
		prToMerge := &domain.PullRequest{
			ID:     prID,
			Status: domain.PRStatusOpen,
		}

		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(prToMerge, nil)
		repo.EXPECT().UpdateStatus(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
				pr.MergedAt = &now
				return pr, nil
			},
		)

		pr, err := uc.MergePullRequest(ctx, prID)
		assert.NoError(t, err)
		assert.Equal(t, domain.PRStatusMerged, pr.Status)
		assert.NotNil(t, pr.MergedAt)
	})
}

func TestReassignReviewer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocksRepo.NewMockPullRequestRepo(ctrl)
	logger := mocksLogger.NewMockLogger(ctrl)
	userRepo := mocksUserRepo.NewMockUserRepo(ctrl)

	uc := PullRequestUsecase{repo: repo, logger: logger, userRepo: userRepo}

	ctx := context.Background()
	prID := 1
	oldReviewer := 10

	pr := &domain.PullRequest{
		ID:                prID,
		AuthorID:          1,
		Status:            domain.PRStatusOpen,
		AssignedReviewers: []int{oldReviewer, 11},
	}

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, oldReviewer).Return(false, nil)

		prResult, newID, err := uc.ReassignReviewer(ctx, &domain.ReassingReviewer{UserID: oldReviewer, PullRequestID: prID})
		assert.Nil(t, prResult)
		assert.Equal(t, 0, newID)
		assert.Equal(t, domain.ErrUserNotFound, err)
	})

	t.Run("PR not found", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, oldReviewer).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, prID).Return(false, nil)

		prResult, newID, err := uc.ReassignReviewer(ctx, &domain.ReassingReviewer{UserID: oldReviewer, PullRequestID: prID})
		assert.Nil(t, prResult)
		assert.Equal(t, 0, newID)
		assert.Equal(t, domain.ErrPullRequestNotFound, err)
	})

	t.Run("PR already merged", func(t *testing.T) {
		mergedPR := &domain.PullRequest{ID: prID, Status: domain.PRStatusMerged, AssignedReviewers: []int{oldReviewer}}
		userRepo.EXPECT().ExistsById(ctx, oldReviewer).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(mergedPR, nil)

		prResult, newID, err := uc.ReassignReviewer(ctx, &domain.ReassingReviewer{UserID: oldReviewer, PullRequestID: prID})
		assert.Nil(t, prResult)
		assert.Equal(t, 0, newID)
		assert.Equal(t, domain.ErrPullRequestIsMerged, err)
	})

	t.Run("user not assigned", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, 99).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(pr, nil)

		prResult, newID, err := uc.ReassignReviewer(ctx, &domain.ReassingReviewer{UserID: 99, PullRequestID: prID})
		assert.Nil(t, prResult)
		assert.Equal(t, 0, newID)
		assert.Equal(t, domain.ErrNotAssigned, err)
	})

	t.Run("no available candidates", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, oldReviewer).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(pr, nil)
		repo.EXPECT().GetActiveTeamMembersExceptAuthor(ctx, pr.AuthorID).Return([]domain.User{
			{ID: oldReviewer},
			{ID: 11},
		}, nil)

		prResult, newID, err := uc.ReassignReviewer(ctx, &domain.ReassingReviewer{UserID: oldReviewer, PullRequestID: prID})
		assert.Nil(t, prResult)
		assert.Equal(t, 0, newID)
		assert.Equal(t, domain.ErrNoAvailableCandidats, err)
	})

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().ExistsById(ctx, oldReviewer).Return(true, nil)
		repo.EXPECT().ExistsById(ctx, prID).Return(true, nil)
		repo.EXPECT().GetById(ctx, prID).Return(pr, nil)
		repo.EXPECT().GetActiveTeamMembersExceptAuthor(ctx, pr.AuthorID).Return([]domain.User{
			{ID: 12}, {ID: 13}, {ID: 14},
		}, nil)
		repo.EXPECT().UpdateAssignedReviewers(ctx, prID, oldReviewer, gomock.Any()).Return(nil)

		prResult, newID, err := uc.ReassignReviewer(ctx, &domain.ReassingReviewer{UserID: oldReviewer, PullRequestID: prID})
		assert.NoError(t, err)
		assert.NotNil(t, prResult)
		assert.NotEqual(t, oldReviewer, newID)
		assert.Contains(t, prResult.AssignedReviewers, newID)
	})
}
