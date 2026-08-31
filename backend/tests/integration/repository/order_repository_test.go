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

type OrderTestData struct {
	User    domain.User
	Product domain.Product
}

func setupOrderTest(t *testing.T, db *sqlx.DB) OrderTestData {
	t.Helper()

	user := helpers.CreateTestUser(t, db, "buyer", "buyer@gmail.com", "1234Qwerty")
	seller := helpers.CreateTestUser(t, db, "seller", "seller@gmail.com", "1234Qwerty")

	category, err := domain.NewCategory("categoryName", "categorySlug")
	require.Nil(t, err)
	categoryID := helpers.CreateTestCategory(t, db, category)

	price, err := decimal.NewFromString("100.50")
	require.Nil(t, err)

	product, err := domain.NewProduct("product", "product description", "image",
		price, categoryID, seller.ID)
	require.Nil(t, err)
	createdProduct := helpers.CreateTestProduct(t, db, product)

	return OrderTestData{
		User:    user,
		Product: createdProduct,
	}
}

func TestOrderRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewOrderRepository(db)
	ctx := context.Background()

	t.Run("Create success", func(t *testing.T) {
		data := setupOrderTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		total, err := decimal.NewFromString("100.50")
		require.Nil(t, err)

		orderID, err := repo.Create(ctx, data.User.ID, config.PendingOrderStatus, total)
		require.Nil(t, err)
		require.Greater(t, orderID, 0)
	})

	t.Run("ListByUserID success", func(t *testing.T) {
		data := setupOrderTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		total, err := decimal.NewFromString("201.00")
		require.Nil(t, err)
		itemTotal, err := decimal.NewFromString("100.50")
		require.Nil(t, err)

		orderID := helpers.CreateTestOrder(t, db, data.User.ID, config.PendingOrderStatus, total)
		helpers.CreateTestOrderItem(t, db, orderID, data.Product, 2, itemTotal.Mul(decimal.NewFromInt(2)))

		orders, err := repo.ListByUserID(ctx, data.User.ID)
		require.Nil(t, err)
		require.Len(t, orders, 1)

		assert.Equal(t, orderID, orders[0].ID)
		assert.Equal(t, data.User.ID, orders[0].UserID)
		assert.Equal(t, config.PendingOrderStatus, orders[0].Status)
		assert.True(t, total.Equal(orders[0].Total))
		require.Len(t, orders[0].OrderItems, 1)
		assert.Equal(t, data.Product.ID, orders[0].OrderItems[0].Product.ID)
		assert.Equal(t, data.Product.Name, orders[0].OrderItems[0].Product.Name)
		assert.Equal(t, 2, orders[0].OrderItems[0].Quantity)
	})

	t.Run("ListByUserID empty", func(t *testing.T) {
		data := setupOrderTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		orders, err := repo.ListByUserID(ctx, data.User.ID)
		require.Nil(t, err)
		assert.Empty(t, orders)
	})
}
