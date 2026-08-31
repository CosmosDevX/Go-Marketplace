package repository

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/tests/helpers"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CartItemTestData struct {
	CartID  int
	Product domain.Product
}

func setupCartItemTest(t *testing.T, db *sqlx.DB) CartItemTestData {
	user := helpers.CreateTestUser(t, db, "user", "user@gmail.com", "1234Qwerty")
	seller := helpers.CreateTestUser(t, db, "seller", "seller@gmail.com", "1234Qwerty")
	cartID := helpers.CreateTestCart(t, db, user.ID)

	category, err := domain.NewCategory("categoryName", "categorySlug")
	require.Nil(t, err)
	categoryID := helpers.CreateTestCategory(t, db, category)

	price, err := decimal.NewFromString("718.91")
	require.Nil(t, err)

	product, err := domain.NewProduct("product", "product description", "image",
		price, categoryID, seller.ID)
	require.Nil(t, err)
	createdProduct := helpers.CreateTestProduct(t, db, product)

	return CartItemTestData{
		CartID:  cartID,
		Product: createdProduct,
	}
}

func TestCartItemRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewCartItemRepository(db)
	ctx := context.Background()

	t.Run("Create and ListByCartID success", func(t *testing.T) {
		data := setupCartItemTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		cartItemID, err := repo.Create(ctx, data.CartID, data.Product.ID)
		require.Nil(t, err)
		require.Greater(t, cartItemID, 0)

		cartItems, err := repo.ListByCartID(ctx, data.CartID)
		require.Nil(t, err)

		assert.Equal(t, len(cartItems), 1)
		assert.Equal(t, cartItems[0].Quantity, 1)
		assert.Equal(t, cartItems[0].Product.ID, data.Product.ID)
	})

	t.Run("Update quantity success", func(t *testing.T) {
		data := setupCartItemTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		cartItemID := helpers.CreateTestCartItem(t, db, data.CartID, data.Product)

		newQuantity, err := repo.UpdateQuantity(ctx, data.CartID, cartItemID, 1)
		require.Nil(t, err)
		assert.Equal(t, newQuantity, 2)

		newQuantity, err = repo.UpdateQuantity(ctx, data.CartID, cartItemID, -1)
		require.Nil(t, err)
		assert.Equal(t, newQuantity, 1)
	})

	t.Run("Delete success", func(t *testing.T) {
		data := setupCartItemTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		cartItemID := helpers.CreateTestCartItem(t, db, data.CartID, data.Product)
		err := repo.Delete(ctx, data.CartID, cartItemID)
		require.Nil(t, err)
	})

	t.Run("Update quantity not found", func(t *testing.T) {
		data := setupCartItemTest(t, db)
		t.Cleanup(func() {
			helpers.TruncateAllTables(t, db)
		})

		newQuantity, err := repo.UpdateQuantity(ctx, data.CartID, 99999, 1)
		require.NotNil(t, err)
		assert.Equal(t, newQuantity, 0)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
