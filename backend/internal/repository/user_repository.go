// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"myapp/internal/constants"
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

func (r UserRepository) GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewDomainError(constants.NotFound, "user not found")
		}

		return nil, domain.NewDomainError(constants.FindError, "error during get user by name")
	}

	return &user, nil
}

func (r UserRepository) CreateUser(ctx context.Context, userDTO dto.CreateUserDTO) (int, *domain.DomainError) {
	query := `INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, userDTO.Username, userDTO.Password, userDTO.Email).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, &domain.DomainError{Code: constants.UniqueViolation, Message: "username or email already taken"}
		}

		return 0, domain.NewDomainError(constants.CreateError, "error during create user")
	}

	return id, nil
}
