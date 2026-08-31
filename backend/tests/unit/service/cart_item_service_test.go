package service_test

import (
	"context"
	"testing"

	"myapp/internal/domain"
	"myapp/internal/service"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCartItemRepo struct {
	createFn         func(ctx context.Context, cartID, productID int) (int, error)
	listByCartIDFn   func(ctx context.Context, cartID int) ([]domain.CartItem, error)
	updateQuantityFn func(ctx context.Context, cartID, cartItemID, delta int) (int, error)
	deleteFn         func(ctx context.Context, cartID, cartItemID int) error
}

func (m *mockCartItemRepo) Create(ctx context.Context, cartID, productID int) (int, error) {
	return m.createFn(ctx, cartID, productID)
}
func (m *mockCartItemRepo) ListByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error) {
	return m.listByCartIDFn(ctx, cartID)
}
func (m *mockCartItemRepo) UpdateQuantity(ctx context.Context, cartID, cartItemID, delta int) (int, error) {
	return m.updateQuantityFn(ctx, cartID, cartItemID, delta)
}
func (m *mockCartItemRepo) Delete(ctx context.Context, cartID, cartItemID int) error {
	return m.deleteFn(ctx, cartID, cartItemID)
}

type mockCartIDGetter struct {
	getIDByUserIDFn func(ctx context.Context, userID int) (int, error)
}

func (m *mockCartIDGetter) GetIDByUserID(ctx context.Context, userID int) (int, error) {
	return m.getIDByUserIDFn(ctx, userID)
}

func TestCartItemService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockCartItemRepo{
			createFn: func(ctx context.Context, cartID, productID int) (int, error) {
				assert.Equal(t, 10, cartID)
				assert.Equal(t, 5, productID)
				return 42, nil
			},
		}
		getter := &mockCartIDGetter{
			getIDByUserIDFn: func(ctx context.Context, userID int) (int, error) {
				assert.Equal(t, 1, userID)
				return 10, nil
			},
		}
		svc := service.NewCartItemService(repo, getter)

		id, err := svc.Create(ctx, service.CreateCartItemInput{UserID: 1, ProductID: 5})
		require.NoError(t, err)
		assert.Equal(t, 42, id)
	})

	t.Run("cart not found", func(t *testing.T) {
		getter := &mockCartIDGetter{
			getIDByUserIDFn: func(ctx context.Context, userID int) (int, error) {
				return 0, domain.ErrNotFound
			},
		}
		svc := service.NewCartItemService(&mockCartItemRepo{}, getter)

		id, err := svc.Create(ctx, service.CreateCartItemInput{UserID: 1, ProductID: 5})
		require.ErrorIs(t, err, domain.ErrNotFound)
		assert.Equal(t, 0, id)
	})
}

func TestCartItemService_List(t *testing.T) {
	ctx := context.Background()
	price := decimal.NewFromInt(100)

	t.Run("success", func(t *testing.T) {
		expected := []domain.CartItem{{
			ID: 1, CartID: 10, Quantity: 2,
			Product: domain.Product{ID: 5, Name: "p", Price: price},
		}}
		repo := &mockCartItemRepo{
			listByCartIDFn: func(ctx context.Context, cartID int) ([]domain.CartItem, error) {
				assert.Equal(t, 10, cartID)
				return expected, nil
			},
		}
		getter := &mockCartIDGetter{
			getIDByUserIDFn: func(ctx context.Context, userID int) (int, error) {
				return 10, nil
			},
		}
		svc := service.NewCartItemService(repo, getter)

		items, err := svc.List(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, expected, items)
	})
}

func TestCartItemService_UpdateQuantity(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockCartItemRepo{
			updateQuantityFn: func(ctx context.Context, cartID, cartItemID, delta int) (int, error) {
				assert.Equal(t, 10, cartID)
				assert.Equal(t, 3, cartItemID)
				assert.Equal(t, 1, delta)
				return 2, nil
			},
		}
		getter := &mockCartIDGetter{
			getIDByUserIDFn: func(ctx context.Context, userID int) (int, error) {
				return 10, nil
			},
		}
		svc := service.NewCartItemService(repo, getter)

		qty, err := svc.UpdateQuantity(ctx, 1, 3, 1)
		require.NoError(t, err)
		assert.Equal(t, 2, qty)
	})

	t.Run("invalid delta", func(t *testing.T) {
		svc := service.NewCartItemService(&mockCartItemRepo{}, &mockCartIDGetter{})

		qty, err := svc.UpdateQuantity(ctx, 1, 3, 5)
		require.ErrorIs(t, err, domain.ErrValidation)
		assert.Equal(t, 0, qty)
	})
}

func TestCartItemService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockCartItemRepo{
			deleteFn: func(ctx context.Context, cartID, cartItemID int) error {
				assert.Equal(t, 10, cartID)
				assert.Equal(t, 3, cartItemID)
				return nil
			},
		}
		getter := &mockCartIDGetter{
			getIDByUserIDFn: func(ctx context.Context, userID int) (int, error) {
				return 10, nil
			},
		}
		svc := service.NewCartItemService(repo, getter)

		err := svc.Delete(ctx, 1, 3)
		require.NoError(t, err)
	})
}
