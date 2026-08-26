package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type cartItemRow struct {
	ID                 int             `db:"cart_item_id"`
	CartID             int             `db:"cart_id"`
	Quantity           int             `db:"quantity"`
	ProductID          int             `db:"product_id"`
	ProductName        string          `db:"product_name"`
	ProductDescription string          `db:"product_description"`
	ProductPrice       decimal.Decimal `db:"product_price"`
	ProductImage       string          `db:"product_image"`
	ProductSellerID    int             `db:"product_seller_id"`
	CategoryID         int             `db:"category_id"`
	CategoryName       string          `db:"category_name"`
	CategorySlug       string          `db:"category_slug"`
}

func (r cartItemRow) toDomain() domain.CartItem {
	return domain.CartItem{
		ID:       r.ID,
		CartID:   r.CartID,
		Quantity: r.Quantity,
		Product: domain.Product{
			ID:          r.ProductID,
			Name:        r.ProductName,
			Description: r.ProductDescription,
			Price:       r.ProductPrice,
			Image:       r.ProductImage,
			SellerID:    r.ProductSellerID,
			Category: domain.Category{
				ID:   r.CategoryID,
				Name: r.CategoryName,
				Slug: r.CategorySlug,
			},
		},
	}
}

type CartItemRepository struct {
	db DBTX
}

func NewCartItemRepository(db DBTX) CartItemRepository {
	return CartItemRepository{
		db: db,
	}
}

func (r CartItemRepository) Create(ctx context.Context, cartID, productID int) (int, error) {
	query := `INSERT INTO cart_items(cart_id, product_id) VALUES($1, $2) RETURNING cart_item_id`
	var id int
	err := r.db.QueryRowContext(ctx, query, cartID, productID).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create cart item: %w", domain.ErrTimeout)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("create cart item: %w", domain.ErrUniqueViolation)
		}

		return 0, fmt.Errorf("create cart item: %w", err)
	}

	return id, nil
}

func (r CartItemRepository) ListByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error) {
	query := `
		SELECT ci.cart_item_id, ci.cart_id, ci.quantity, 
		p.product_id, p.product_name, p.product_description, p.product_price, p.product_image, p.product_seller_id,
		c.category_id, c.category_name, c.category_slug
		FROM cart_items AS ci
		JOIN products AS p ON p.product_id = ci.product_id
		JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE ci.cart_id = $1
	`
	var cartItemRows []cartItemRow
	err := r.db.SelectContext(ctx, &cartItemRows, query, cartID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("get cart items by cartID: %w", domain.ErrTimeout)
		}

		return nil, fmt.Errorf("get cart items by cartID: %w", err)
	}

	domainModels := make([]domain.CartItem, len(cartItemRows))
	for i := range domainModels {
		domainModels[i] = cartItemRows[i].toDomain()
	}

	return domainModels, nil
}

func (r CartItemRepository) UpdateQuantity(ctx context.Context, cartID, cartItemID, delta int) (int, error) {
	query := `UPDATE cart_items SET quantity = quantity + $1 WHERE cart_id = $2 AND cart_item_id = $3 RETURNING quantity`
	var quantity int
	err := r.db.QueryRowContext(ctx, query, delta, cartID, cartItemID).Scan(&quantity)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("update quantity by %d on cart item %d: %w", delta, cartItemID, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("update quantity by %d on cart item %d: %w", delta, cartItemID, domain.ErrNotFound)
		}

		return 0, fmt.Errorf("update quantity by %d on cart item %d: %w", delta, cartItemID, err)
	}

	if quantity <= 0 {
		if err := r.Delete(ctx, cartID, cartItemID); err != nil {
			return 0, err
		}

		return quantity, nil
	}

	return quantity, nil
}

func (r CartItemRepository) Delete(ctx context.Context, cartID, cartItemID int) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1 AND cart_item_id = $2`
	sqlResult, err := r.db.ExecContext(ctx, query, cartID, cartItemID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("delete cart item by id %d in cart %d: %w", cartItemID, cartID, err)
		}

		return fmt.Errorf("delete cart item by id %d in cart %d: %w", cartItemID, cartID, err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("cart item %d in cart %d: %w", cartItemID, cartID, domain.ErrNotFound)
	}

	return nil
}

func (r CartItemRepository) DeleteAllByCartID(ctx context.Context, cartID int) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	sqlResult, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("clear cart %d: %w", cartID, domain.ErrTimeout)
		}

		return fmt.Errorf("clear cart %d: %w", cartID, err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("clear cart %d: %w", cartID, domain.ErrNotFound)
	}

	return nil
}
