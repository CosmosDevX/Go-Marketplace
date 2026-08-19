// Package repository
package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
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

func (r UserRoleRepository) Create(ctx context.Context, userID int, roleName string) error {
	query := `
		INSERT INTO user_roles(role_id, user_id) 
		SELECT r.role_id, $1
		FROM roles AS r
		WHERE r.role_name = $2
	`
	sqlResult, err := r.db.ExecContext(ctx, query, userID, roleName)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("create user role: %w", domain.ErrTimeout)
		}

		return fmt.Errorf("create user role: %w", err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("role %s: %w", roleName, domain.ErrNotFound)
	}

	return nil
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
