package user

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

const (
	checkUserById = `
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);
	`

	updateUserIsActive = `
		UPDATE users u SET is_active = $1
		FROM team t 
		WHERE u.id = $2 AND u.team_id = t.id
		RETURNING u.id, u.name, u.is_active, t.name;
	`

	getUserPullRequests = `
		SELECT pr.id, pr.title, pr.author_id, s.name, pr.created_at, pr.merged_at
		FROM pull_request pr
		JOIN pr_status s ON pr.status_id = s.id
		WHERE pr.id IN (
			SELECT pr_id
			FROM assigned_pr
			WHERE reviewer_id = $1
		)
		ORDER BY pr.created_at DESC;
	`
)

func (r *UserRepository) ExistsById(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, checkUserById, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existance by id: %w", err)
	}
	return exists, nil
}

func (r *UserRepository) UpdateIsActive(ctx context.Context, set *domain.SetUserIsActive) (*domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, updateUserIsActive, set.IsActive, set.ID).
		Scan(&user.ID, &user.Username, &user.IsActive, &user.TeamName)

	if err != nil {
		return nil, fmt.Errorf("failed to update user is_active: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserPullRequests(ctx context.Context, userID int) ([]domain.PullRequest, error) {
	rows, err := r.pool.Query(ctx, getUserPullRequests, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user pull_requests: %w", err)
	}
	defer rows.Close()

	var prs []domain.PullRequest
	for rows.Next() {
		var status string
		var pr domain.PullRequest

		err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &status, &pr.CreatedAt, &pr.MergedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pull_request: %w", err)
		}
		pr.Status = domain.MapStringToPullRequestStatus[status]

		prs = append(prs, pr)
	}

	return prs, nil
}
