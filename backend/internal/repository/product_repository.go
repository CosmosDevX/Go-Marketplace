package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) ProductRepository {
	return ProductRepository{
		db: db,
	}
}

type CreateProductParams struct {
	Name        string
	Description string
	Price       decimal.Decimal
	Image       string
	CategoryID  int
}

const productPageSize = 16

func (r ProductRepository) ListAll(ctx context.Context, page int) ([]domain.Product, error) {
	query := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price, p.product_image,
			c.category_id AS "category.category_id", c.category_name AS "category.category_name", c.category_slug AS "category.category_slug"
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		LIMIT 16 OFFSET $1
	`
	offset := productPageSize * (page - 1)

	var products []domain.Product
	err := r.db.SelectContext(ctx, &products, query, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get all products: %w", domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get all products: %w", err)

	}

	return products, nil
}

func (r ProductRepository) ListByCategorySlug(ctx context.Context, categorySlug string, page int) ([]domain.Product, error) {
	query := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price, p.product_image,
			c.category_id AS "category.category_id", c.category_name AS "category.category_name", c.category_slug AS "category.category_slug"
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE c.category_slug = $1
		LIMIT 16 OFFSET $2
	`
	offset := productPageSize * (page - 1)

	var products []domain.Product
	err := r.db.SelectContext(ctx, &products, query, categorySlug, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get products by category slug %s: %w", categorySlug, domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get products by category slug %s: %w", categorySlug, err)
	}

	return products, nil
}

func (r ProductRepository) Create(ctx context.Context, params CreateProductParams) (int, error) {
	query := `
		INSERT INTO products(product_name, product_description, product_price, product_image, product_category_id)
	 	VALUES($1, $2, $3, $4, $5) RETURNING product_id
	`
	var productID int
	err := r.db.QueryRowContext(ctx, query, params.Name, params.Description, params.Price, params.Image, params.CategoryID).Scan(&productID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create product: %w", domain.ErrTimeout)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("create product: %w", domain.ErrUniqueViolation)
		}

		return 0, fmt.Errorf("create product: %w", err)
	}

	return productID, nil
}
