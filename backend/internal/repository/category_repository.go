// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/lib/pq"
)

type CategoryRepository struct {
	db DBTX
}

func NewCategoryRepository(db DBTX) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

type CreateCategoryParams struct {
	Name string
	Slug string
}

func (r CategoryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT category_id, category_name, category_slug FROM categories`
	var categories []domain.Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get all categories: %w", domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get all categories: %w", err)
	}

	return categories, nil
}

func (r CategoryRepository) Create(ctx context.Context, params CreateCategoryParams) (int, error) {
	query := `INSERT INTO categories(category_name, category_slug) VALUES($1, $2) RETURNING category_id`
	var categoryID int
	err := r.db.QueryRowContext(ctx, query, params.Name, params.Slug).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create category: %w", domain.ErrTimeout)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("create category: %w", domain.ErrUniqueViolation)
		}

		return 0, fmt.Errorf("create category: %w", err)
	}

	return categoryID, nil
}

func (r CategoryRepository) GetIDBySlug(ctx context.Context, categorySlug string) (int, error) {
	query := `SELECT category_id FROM categories WHERE category_slug = $1`
	var categoryID int
	err := r.db.QueryRowContext(ctx, query, categorySlug).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("get category id by slug %s: %w", categorySlug, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("get category id by slug %s: %w", categorySlug, domain.ErrNotFound)
		}

		return 0, fmt.Errorf("get category id by slug %s: %w", categorySlug, err)
	}

	return categoryID, nil
}
