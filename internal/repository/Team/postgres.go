package team

import (
	"context"
	"fmt"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamPepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewTeamRepository(pool *pgxpool.Pool, logger logger.Logger) *TeamPepository {
	return &TeamPepository{
		pool:   pool,
		logger: logger,
	}
}

const (
	checkTeamByName = `
		SELECT EXISTS(SELECT 1 FROM team WHERE name = $1);
	`
	checkUserByID = `
        SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);
    `

	createTeamWithName = `
		INSERT INTO team (name) VALUES ($1) RETURNING id;
	`

	createTeamMember = `
		INSERT INTO users (id, name, is_active, team_id) VALUES ($1, $2, $3, $4);
	`

	updateTeamMember = `
		UPDATE users SET name = $1, is_active = $2, team_id = $3 WHERE id = $4;
	`

	getTeamByName = `
		SELECT id, name FROM team WHERE name = $1;
	`

	getTeamMembers = `
		SELECT id, name, is_active FROM users WHERE team_id = $1;
	`
)

// Проверка существования команды с заданным именем
func (r *TeamPepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, checkTeamByName, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Создание команды с созданием/обновлением участников
func (r *TeamPepository) Create(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			r.logger.WithFields(logger.LoggerFields{"err": err.Error()}).Error("tx rollback failed")
		}
	}()

	// Создаем команду
	var teamID int
	err = tx.QueryRow(ctx, createTeamWithName, team.Name).Scan(&teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert team: %w", err)
	}

	for _, m := range team.Members {
		// Проверяем существование User с id
		var exists bool
		err := tx.QueryRow(ctx, checkUserByID, m.UserID).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check User by id: %w", err)
		}

		// Если существует - обновляем
		if exists {
			_, err := tx.Exec(ctx, updateTeamMember, m.Username, m.IsActive, teamID)
			if err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
			// Если не существует - создаем
		} else {
			_, err := tx.Exec(ctx, createTeamMember, m.UserID, m.Username, m.IsActive, teamID)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	team.ID = teamID

	return team, nil
}

func (r *TeamPepository) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	var team domain.Team

	err := r.pool.QueryRow(ctx, getTeamByName, name).Scan(&team.ID, &team.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	rows, err := r.pool.Query(ctx, getTeamMembers, team.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.TeamMember
		err := rows.Scan(&m.UserID, &m.Username, &m.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		team.Members = append(team.Members, m)
	}

	return &team, nil
}
