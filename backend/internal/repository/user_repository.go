// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (r UserRepository) GetByName(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get user by name %s: %w", username, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by name %s: %w", username, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("get user by name %s: %w", username, err)
	}

	return &user, nil
}

func (r UserRepository) Create(ctx context.Context, userDTO dto.CreateUserDTO) (int, error) {
	query := `INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, userDTO.Username, userDTO.Password, userDTO.Email).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create user: %w", domain.ErrTimeout)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("create user: %w", domain.ErrUniqueViolation)
		}

		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}
