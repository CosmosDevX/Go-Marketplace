// Package repository
package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/lib/pq"
)

type userRoleRow struct {
	ID     int `db:"user_role_id"`
	RoleID int `db:"role_id"`
	UserID int `db:"user_id"`
}

type UserRoleRepository struct {
	db DBTX
}

func NewUserRoleRepository(db DBTX) UserRoleRepository {
	return UserRoleRepository{
		db: db,
	}
}

func (r UserRoleRepository) Create(ctx context.Context, username, roleName string) error {
	query := `
		INSERT INTO user_roles(role_id, user_id) 
		SELECT r.role_id, (SELECT u.id FROM users AS u WHERE u.username = $1)
		FROM roles AS r
		WHERE r.role_name = $2
	`
	sqlResult, err := r.db.ExecContext(ctx, query, username, roleName)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("create user role: %w", domain.ErrTimeout)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("create user role: %w", domain.ErrUniqueViolation)
		}

		return fmt.Errorf("create user role: %w", err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("role %s: %w", roleName, domain.ErrNotFound)
	}

	return nil
}

func (r UserRoleRepository) Delete(ctx context.Context, username, roleName string) error {
	query := `
		DELETE FROM user_roles AS ur
		WHERE ur.user_id = (SELECT u.id FROM users AS u WHERE u.username = $1)
		AND ur.role_id = (SELECT r.role_id FROM roles AS r WHERE r.role_name = $2) 	
	`
	sqlResult, err := r.db.ExecContext(ctx, query, username, roleName)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("delete role %s from username %s: %w", roleName, username, err)
		}

		return fmt.Errorf("delete role %s from username %s: %w", roleName, username, err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("role %s from username %s: %w", roleName, username, domain.ErrNotFound)
	}

	return nil
}

func (r UserRoleRepository) ListByUsername(ctx context.Context, username string) ([]string, error) {
	query := `
		SELECT r.role_name
		FROM user_roles ur
		JOIN roles r ON r.role_id = ur.role_id
		WHERE ur.user_id = (SELECT u.id FROM users AS u WHERE u.username = $1)
	`
	var roles []string
	if err := r.db.SelectContext(ctx, &roles, query, username); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("list roles by username %s: %w", username, domain.ErrTimeout)
		}

		return nil, fmt.Errorf("list roles by username %s: %w", username, err)
	}

	return roles, nil
}

func (r UserRoleRepository) ListByUserID(ctx context.Context, userID int) ([]string, error) {
	query := `
		SELECT r.role_name
		FROM user_roles ur
		JOIN roles r ON r.role_id = ur.role_id
		WHERE ur.user_id = $1
	`
	var roles []string
	if err := r.db.SelectContext(ctx, &roles, query, userID); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("list roles by userID %d: %w", userID, domain.ErrTimeout)
		}

		return nil, fmt.Errorf("list roles by userID %d: %w", userID, err)
	}

	return roles, nil
}
