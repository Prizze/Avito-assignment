package pullrequest

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{
		pool: pool,
	}
}

const (
	checkRPById = `
		SELECT EXISTS (SELECT 1 FROM pull_request WHERE id = $1);
	`

	getActiveTeamMembers = `
		SELECT id, name, is_active, (SELECT name from team WHERE id = users.team_id)
		FROM users
		WHERE team_id = (SELECT team_id FROM users WHERE id = $1)
			AND id <> $1
			AND is_active = TRUE; 
	`

	createPullRequest = `
		INSERT INTO pull_request (id, title, author_id, status_id, created_at)
		VALUES($1, $2, $3, $4, $5);
	`

	getStatusID = `
		SELECT id FROM pr_status WHERE name = $1;
	`

	addReviewerToPullRequest = `
		INSERT INTO assigned_pr (pr_id, reviewer_id) VALUES ($1, $2);
	`

	getPullRequestByID = `
		SELECT id, title, author_id, 
		(SELECT name FROM pr_status WHERE pr_status.id = pull_request.status_id) as status,
		created_at, merged_at
		FROM pull_request WHERE id = $1;
	`

	getReviewers = `
		SELECT reviewer_id FROM assigned_pr WHERE pr_id = $1;
	`

	updateStatus = `
		UPDATE pull_request SET status_id = (SELECT id FROM pr_status WHERE name = $1),
		merged_at = $2
		WHERE id = $3;
	`
)

func (r *PullRequestRepository) ExistsById(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, checkRPById, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PullRequestRepository) GetActiveTeamMembersExceptAuthor(ctx context.Context, authorId int) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx, getActiveTeamMembers, authorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get active team members: %w", err)
	}
	defer rows.Close()

	activeMembers := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.IsActive,
			&u.TeamName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		activeMembers = append(activeMembers, u)
	}

	return activeMembers, nil
}

func (r *PullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var statusID int
	err = tx.QueryRow(ctx, getStatusID, pr.Status).Scan(&statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status_id: %w", err)
	}

	_, err = tx.Exec(ctx, createPullRequest, pr.ID, pr.Name, pr.AuthorID, statusID, pr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert pull_request: %w", err)
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err := tx.Exec(ctx, addReviewerToPullRequest, pr.ID, reviewerID)
		if err != nil {
			return nil, fmt.Errorf("faield to insert reviewer: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return pr, nil
}

func (r *PullRequestRepository) GetById(ctx context.Context, id int) (*domain.PullRequest, error) {
	pr := &domain.PullRequest{}
	var status string

	err := r.pool.QueryRow(ctx, getPullRequestByID, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &status, &pr.CreatedAt, &pr.MergedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull_request: %w", err)
	}

	pr.Status = domain.MapStringToPullRequestStatus[status]

	rows, err := r.pool.Query(ctx, getReviewers, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reviewerID int
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	}

	return pr, nil
}

func (r *PullRequestRepository) UpdateStatus(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	_, err := r.pool.Exec(ctx, updateStatus, pr.Status, pr.MergedAt, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update pull_request status: %w", err)
	}

	return pr, nil
}
