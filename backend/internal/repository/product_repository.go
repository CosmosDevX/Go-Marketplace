package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type productRow struct {
	ID           int             `db:"product_id"`
	Name         string          `db:"product_name"`
	Description  string          `db:"product_description"`
	Price        decimal.Decimal `db:"product_price"`
	Image        string          `db:"product_image"`
	CategoryID   int             `db:"category_id"`
	CategoryName string          `db:"category_name"`
	CategorySlug string          `db:"category_slug"`
}

func (r productRow) toDomain() domain.Product {
	return domain.Product{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Image:       r.Image,
		Category: domain.Category{
			ID:   r.CategoryID,
			Name: r.CategoryName,
			Slug: r.CategorySlug,
		},
	}
}

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db DBTX) ProductRepository {
	return ProductRepository{
		db: db,
	}
}

func (r ProductRepository) List(ctx context.Context, page int) ([]domain.Product, error) {
	query := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price, p.product_image,
			c.category_id, c.category_name, c.category_slug
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		LIMIT $1 OFFSET $2
	`
	offset := config.ProductPageSize * (page - 1)

	var productRows []productRow
	err := r.db.SelectContext(ctx, &productRows, query, config.ProductPageSize, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get all products: %w", domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get all products: %w", err)
	}

	domainModels := make([]domain.Product, len(productRows))
	for i := range domainModels {
		domainModels[i] = productRows[i].toDomain()
	}

	return domainModels, nil
}

func (r ProductRepository) ListByCategorySlug(ctx context.Context, categorySlug string, page int) ([]domain.Product, error) {
	query := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price, p.product_image,
			c.category_id, c.category_name, c.category_slug
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE c.category_slug = $1
		LIMIT $2 OFFSET $3
	`
	offset := config.ProductPageSize * (page - 1)
	var productRows []productRow
	err := r.db.SelectContext(ctx, &productRows, query, categorySlug, config.ProductPageSize, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get products by category slug %s: %w", categorySlug, domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get products by category slug %s: %w", categorySlug, err)
	}

	domainModels := make([]domain.Product, len(productRows))
	for i := range domainModels {
		domainModels[i] = productRows[i].toDomain()
	}

	return domainModels, nil
}

func (r ProductRepository) Create(ctx context.Context, p domain.Product) (int, error) {
	query := `
		INSERT INTO products(product_name, product_description, product_price, product_image, product_category_id)
	 	VALUES($1, $2, $3, $4, $5) RETURNING product_id
	`
	var productID int
	err := r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price, p.Image, p.Category.ID).Scan(&productID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create product: %w", domain.ErrTimeout)
		}

		return 0, fmt.Errorf("create product: %w", err)
	}

	return productID, nil
}
