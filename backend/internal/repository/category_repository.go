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

type categoryRow struct {
	ID   int    `db:"category_id"`
	Name string `db:"category_name"`
	Slug string `db:"category_slug"`
}

func (r categoryRow) toDomain() domain.Category {
	return domain.Category{
		ID:   r.ID,
		Name: r.Name,
		Slug: r.Slug,
	}
}

type CategoryRepository struct {
	db DBTX
}

func NewCategoryRepository(db DBTX) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

func (r CategoryRepository) Create(ctx context.Context, c domain.Category) (int, error) {
	query := `INSERT INTO categories(category_name, category_slug) VALUES($1, $2) RETURNING category_id`
	var categoryID int
	err := r.db.QueryRowContext(ctx, query, c.Name, c.Slug).Scan(&categoryID)
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

func (r CategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT category_id, category_name, category_slug FROM categories`
	var categoryRows []categoryRow
	err := r.db.SelectContext(ctx, &categoryRows, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get all categories: %w", domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get all categories: %w", err)
	}

	domainModels := make([]domain.Category, len(categoryRows))
	for i := range domainModels {
		domainModels[i] = categoryRows[i].toDomain()
	}

	return domainModels, nil
}

func (r CategoryRepository) GetIDBySlug(ctx context.Context, categorySlug string) (int, error) {
	query := `SELECT category_id FROM categories WHERE category_slug = $1`
	var categoryID int
	err := r.db.QueryRowContext(ctx, query, categorySlug).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("get categoryID by slug %s: %w", categorySlug, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("get categoryID by slug %s: %w", categorySlug, domain.ErrNotFound)
		}

		return 0, fmt.Errorf("get categoryID by slug %s: %w", categorySlug, err)
	}

	return categoryID, nil
}
