package repository

import (
	"context"
	"myapp/internal/config"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/tests/helpers"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type OrderItemTestData struct {
	OrderID int
	Product domain.Product
}

func setupOrderItemTest(t *testing.T, db *sqlx.DB) OrderItemTestData {
	t.Helper()

	user := helpers.CreateTestUser(t, db, "buyer", "buyer@gmail.com", "1234Qwerty")
	seller := helpers.CreateTestUser(t, db, "seller", "seller@gmail.com", "1234Qwerty")

	category, err := domain.NewCategory("categoryName", "categorySlug")
	require.Nil(t, err)
	categoryID := helpers.CreateTestCategory(t, db, category)

	price, err := decimal.NewFromString("50.00")
	require.Nil(t, err)

	product, err := domain.NewProduct("product", "product description", "image",
		price, categoryID, seller.ID)
	require.Nil(t, err)
	createdProduct := helpers.CreateTestProduct(t, db, product)

	total, err := decimal.NewFromString("100.00")
	require.Nil(t, err)

	var orderID int
	err = db.QueryRowContext(context.Background(), `
		INSERT INTO orders(user_id, order_status_id, order_total)
		SELECT $1, (SELECT order_status_id FROM order_statuses WHERE order_status_name = $2), $3
		RETURNING order_id
	`, user.ID, config.PendingOrderStatus, total).Scan(&orderID)
	require.NoError(t, err)

	return OrderItemTestData{
		OrderID: orderID,
		Product: createdProduct,
	}
}

func TestOrderItemRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewOrderItemRepository(db)
	ctx := context.Background()

	t.Run("CreateMany success", func(t *testing.T) {
		data := setupOrderItemTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		itemTotal := data.Product.Price.Mul(decimal.NewFromInt(2))
		orderItems := []domain.OrderItem{
			{
				OrderID:  data.OrderID,
				Product:  data.Product,
				Quantity: 2,
				Total:    itemTotal,
			},
		}

		err := repo.CreateMany(ctx, orderItems)
		require.Nil(t, err)

		var count int
		err = db.GetContext(ctx, &count, `SELECT COUNT(*) FROM order_items WHERE order_id = $1`, data.OrderID)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		var quantity int
		var storedTotal decimal.Decimal
		err = db.QueryRowContext(ctx, `
			SELECT order_item_quantity, order_item_total
			FROM order_items WHERE order_id = $1
		`, data.OrderID).Scan(&quantity, &storedTotal)
		require.NoError(t, err)
		assert.Equal(t, 2, quantity)
		assert.True(t, itemTotal.Equal(storedTotal))
	})

	t.Run("CreateMany multiple items", func(t *testing.T) {
		data := setupOrderItemTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		seller2 := helpers.CreateTestUser(t, db, "seller2", "seller2@gmail.com", "1234Qwerty")
		price2, err := decimal.NewFromString("30.00")
		require.Nil(t, err)
		product2, err := domain.NewProduct("product2", "desc2", "image2",
			price2, data.Product.Category.ID, seller2.ID)
		require.Nil(t, err)
		createdProduct2 := helpers.CreateTestProduct(t, db, product2)

		orderItems := []domain.OrderItem{
			{
				OrderID:  data.OrderID,
				Product:  data.Product,
				Quantity: 1,
				Total:    data.Product.Price,
			},
			{
				OrderID:  data.OrderID,
				Product:  createdProduct2,
				Quantity: 3,
				Total:    price2.Mul(decimal.NewFromInt(3)),
			},
		}

		err = repo.CreateMany(ctx, orderItems)
		require.Nil(t, err)

		var count int
		err = db.GetContext(ctx, &count, `SELECT COUNT(*) FROM order_items WHERE order_id = $1`, data.OrderID)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("CreateMany empty validation error", func(t *testing.T) {
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		err := repo.CreateMany(ctx, []domain.OrderItem{})
		require.NotNil(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}
