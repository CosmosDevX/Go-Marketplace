package repository

import (
	"context"
	"database/sql"
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
	SellerID     int             `db:"product_seller_id"`
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
		SellerID: r.SellerID,
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

func (r ProductRepository) List(ctx context.Context, search, categorySlug, sortBy string, asc bool, page int) ([]domain.Product, error) {
	baseQuery := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price, p.product_image, p.product_seller_id,
			c.category_id, c.category_name, c.category_slug
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE %s
		ORDER BY %s
		LIMIT $1 OFFSET $2
	` //TODO: sql injection in where!
	offset := config.ProductPageSize * (page - 1)

	orderBy := r.buildOrderBy(sortBy, asc)
	where := r.buildWhere(search, categorySlug)
	query := fmt.Sprintf(baseQuery, where, orderBy)

	var productRows []productRow
	err := r.db.SelectContext(ctx, &productRows, query, config.ProductPageSize, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get all products: %w", domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get all products: %w", err)
	}

	return r.toDomainModels(productRows), nil
}

func (r ProductRepository) buildOrderBy(sortBy string, asc bool) string {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}

	switch sortBy {
	case "price":
		return fmt.Sprintf("p.product_price %s, p.product_id ASC", dir)
	case "name":
		return fmt.Sprintf("p.product_name %s, p.product_id ASC", dir)
	default:
		return "p.product_name ASC, p.product_id ASC"
	}
}

func (r ProductRepository) buildWhere(search, categorySlug string) string {
	var where string
	if categorySlug == "" {
		where = "1 = 1"
	} else {
		where = fmt.Sprintf("c.category_slug = '%s'", categorySlug)
	}
	if search != "" {
		where += fmt.Sprintf(" AND p.product_name LIKE '%%%s%%'", search)
	}

	return where
}

func (r ProductRepository) ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error) {
	query := `
		SELECT 
			p.product_id, p.product_name, p.product_description, p.product_price, p.product_image, p.product_seller_id,
			c.category_id, c.category_name, c.category_slug
		FROM products AS p
		JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE p.product_seller_id = $1
		LIMIT $2 OFFSET $3
	`
	offset := config.ProductPageSize * (page - 1)
	var productRows []productRow
	err := r.db.SelectContext(ctx, &productRows, query, sellerID, config.ProductPageSize, offset)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get products by sellerID %d: %w", sellerID, domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get products by sellerID %d: %w", sellerID, err)
	}

	return r.toDomainModels(productRows), nil
}

func (r ProductRepository) Create(ctx context.Context, p domain.Product) (int, error) {
	query := `
		INSERT INTO products(product_name, product_description, product_price, product_image, product_category_id, product_seller_id)
	 	VALUES($1, $2, $3, $4, $5, $6) RETURNING product_id
	`
	var productID int
	err := r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price, p.Image, p.Category.ID, p.SellerID).Scan(&productID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create product: %w", domain.ErrTimeout)
		}

		return 0, fmt.Errorf("create product: %w", err)
	}

	return productID, nil
}

func (r ProductRepository) Delete(ctx context.Context, productID, sellerID int, isAdmin bool) error {
	var query string
	var args []any
	if isAdmin {
		query = `DELETE FROM products WHERE product_id = $1`
		args = []any{productID}
	} else {
		query = `DELETE FROM products WHERE product_id = $1 AND product_seller_id = $2`
		args = []any{productID, sellerID}
	}

	sqlResult, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("delete product by id %d: %w", productID, domain.ErrTimeout)
		}

		return fmt.Errorf("delete product by id %d: %w", productID, err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("delete product by id %d: %w", productID, domain.ErrNotFound)
	}

	return nil
}

func (r ProductRepository) GetImageByID(ctx context.Context, productID, sellerID int, isAdmin bool) (string, error) {
	var query string
	var args []any
	if isAdmin {
		query = `SELECT product_image FROM products WHERE product_id = $1`
		args = []any{productID}
	} else {
		query = `SELECT product_image FROM products WHERE product_id = $1 AND product_seller_id = $2`
		args = []any{productID, sellerID}
	}

	var image string
	if err := r.db.GetContext(ctx, &image, query, args...); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("get image by productID %d: %w", productID, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("get image by productID %d: %w", productID, domain.ErrNotFound)
		}

		return "", fmt.Errorf("get image by productID %d: %w", productID, err)
	}

	return image, nil
}

func (r ProductRepository) toDomainModels(rows []productRow) []domain.Product {
	domainModels := make([]domain.Product, len(rows))
	for i := range domainModels {
		domainModels[i] = rows[i].toDomain()
	}

	return domainModels
}
