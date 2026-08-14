package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/domain"

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

func (r ProductRepository) GetProductsByCategory(ctx context.Context, categorySlug string, page int) ([]domain.Product, *domain.DomainError) {
	if page <= 0 {
		return nil, domain.NewDomainError(constants.FindError, "page cannot be negative")
	}

	query := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price,
			c.category_id AS "category.category_id", c.category_name AS "category.category_name", c.category_slug AS "category.category_slug"
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE c.category_slug = $1
		LIMIT 16 OFFSET $2
	`
	var offset int
	if page == 1 {
		offset = 0
	} else {
		offset = 16 * (page - 1)
	}

	var products []domain.Product
	err := r.db.SelectContext(ctx, &products, query, categorySlug, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}

		return nil, domain.NewDomainError(constants.FindError, "error during get all products")
	}

	return products, nil
}

func (r ProductRepository) CreateProduct(ctx context.Context, params CreateProductParams) (int, *domain.DomainError) {
	query := `
		INSERT INTO products(product_name, product_description, product_price, product_image, product_category_id)
	 	VALUES($1, $2, $3, $4, $5) RETURNING product_id
	`
	var productID int
	err := r.db.QueryRowContext(ctx, query, params.Name, params.Description, params.Price, params.Image, params.CategoryID).Scan(&productID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, domain.NewDomainError(constants.RequestTimeout, "request timeout")
		}

		return 0, domain.NewDomainError(constants.CreateError, "error during create product")
	}

	return productID, nil
}
