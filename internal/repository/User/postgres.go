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
