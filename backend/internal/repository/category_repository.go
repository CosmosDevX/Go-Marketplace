// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/domain"
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

func (r CategoryRepository) GetAllCategories(ctx context.Context) ([]domain.Category, *domain.DomainError) {
	query := `SELECT category_id, category_name, category_slug FROM categories`
	var categories []domain.Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}

		return nil, domain.NewDomainError(constants.FindError, "error during get all categories")
	}

	return categories, nil
}

func (r CategoryRepository) CreateCategory(ctx context.Context, params CreateCategoryParams) (int, *domain.DomainError) {
	query := `INSERT INTO categories(category_name, category_slug) VALUES($1, $2) RETURNING category_id`
	var categoryID int
	err := r.db.QueryRowContext(ctx, query, params.Name, params.Slug).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}

		return 0, domain.NewDomainError(constants.CreateError, "error during create category")
	}

	return categoryID, nil
}

func (r CategoryRepository) GetCategoryIDBySlug(ctx context.Context, categorySlug string) (int, *domain.DomainError) {
	query := `SELECT category_id FROM categories WHERE category_slug = $1`
	var categoryID int
	err := r.db.QueryRowContext(ctx, query, categorySlug).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.NewDomainError(constants.NotFound, "category by id not found")
		}

		return 0, domain.NewDomainError(constants.FindError, "error during get category id by slug")
	}

	return categoryID, nil
}
