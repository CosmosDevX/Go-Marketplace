// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRow struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

func (r userRow) toDomain() domain.User {
	return domain.User{
		ID:           r.ID,
		Username:     r.Username,
		PasswordHash: r.Password,
		Email:        r.Email,
	}
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (r UserRepository) GetByName(ctx context.Context, username string) (domain.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	var userRow userRow
	err := r.db.GetContext(ctx, &userRow, query, username)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return domain.User{}, fmt.Errorf("get user by name %s: %w", username, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("get user by name %s: %w", username, domain.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("get user by name %s: %w", username, err)
	}

	domainModel := userRow.toDomain()
	return domainModel, nil
}

func (r UserRepository) Create(ctx context.Context, u domain.User) (int, error) {
	query := `INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, u.Username, u.PasswordHash, u.Email).Scan(&id)
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
