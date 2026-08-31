// Package helpers
package helpers

import (
	"context"
	"testing"

	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func CreateTestUser(t *testing.T, db *sqlx.DB, username, email, password string) domain.User {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	require.NoError(t, err)

	var id int
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`,
		username, email, string(hashed),
	).Scan(&id)
	require.NoError(t, err)

	return domain.User{
		ID:           id,
		Username:     username,
		PasswordHash: string(hashed),
		Email:        email,
	}
}

func CreateTestCart(t *testing.T, db *sqlx.DB, userID int) int {
	t.Helper()

	var id int
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO carts(user_id) VALUES($1) RETURNING cart_id`,
		userID).Scan(&id)
	require.NoError(t, err)

	return id
}

func CreateTestProduct(t *testing.T, db *sqlx.DB, p domain.Product) domain.Product {
	t.Helper()

	var id int
	err := db.QueryRowContext(context.Background(), `
		INSERT INTO products(product_name, product_description, product_price, product_image, product_category_id, product_seller_id)
	 	VALUES($1, $2, $3, $4, $5, $6) RETURNING product_id`,
		p.Name, p.Description, p.Price, p.Image, p.Category.ID, p.SellerID).Scan(&id)
	require.NoError(t, err)
	p.ID = id

	return p
}

func CreateTestCartItem(t *testing.T, db *sqlx.DB, cartID int, p domain.Product) int {
	t.Helper()
	var productID int
	err := db.QueryRowContext(context.Background(), `
		INSERT INTO products(product_name, product_description, product_price, product_image, product_category_id, product_seller_id)
	 	VALUES($1, $2, $3, $4, $5, $6) RETURNING product_id`,
		p.Name, p.Description, p.Price, p.Image, p.Category.ID, p.SellerID).Scan(&productID)
	require.NoError(t, err)

	var cartItemID int
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO cart_items(cart_id, product_id) VALUES($1, $2) RETURNING cart_item_id`,
		cartID, productID).Scan(&cartItemID)
	require.NoError(t, err)

	return cartItemID
}

func CreateTestCategory(t *testing.T, db *sqlx.DB, c domain.Category) int {
	t.Helper()

	var id int
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO categories(category_name, category_slug) VALUES($1, $2) RETURNING category_id`,
		c.Name, c.Slug).Scan(&id)
	require.NoError(t, err)

	return id
}

func CreateTestOrder(t *testing.T, db *sqlx.DB, userID int, status string, total decimal.Decimal) int {
	t.Helper()

	var id int
	err := db.QueryRowContext(context.Background(), `
		INSERT INTO orders(user_id, order_status_id, order_total)
		SELECT $1, (SELECT order_status_id FROM order_statuses WHERE order_status_name = $2), $3
		RETURNING order_id
	`, userID, status, total).Scan(&id)
	require.NoError(t, err)
	return id
}

func CreateTestOrderItem(t *testing.T, db *sqlx.DB, orderID int, product domain.Product, quantity int, total decimal.Decimal) int {
	t.Helper()

	var id int
	err := db.QueryRowContext(context.Background(), `
		INSERT INTO order_items(order_id, product_id, order_item_total, order_item_quantity, product_name, product_price, product_seller_id)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING order_item_id
	`, orderID, product.ID, total, quantity, product.Name, product.Price, product.SellerID).Scan(&id)
	require.NoError(t, err)
	return id
}
