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

type mockUnitOfWork struct {
	doFn func(ctx context.Context, fn func(ctx context.Context, repos service.Repositories) (any, error)) (any, error)
}

func (m *mockUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos service.Repositories) (any, error)) (any, error) {
	return m.doFn(ctx, fn)
}

type mockOrdersGetter struct {
	listByUserIDFn func(ctx context.Context, userID int) ([]domain.Order, error)
}

func (m *mockOrdersGetter) ListByUserID(ctx context.Context, userID int) ([]domain.Order, error) {
	return m.listByUserIDFn(ctx, userID)
}

func TestOrderService_ListByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []domain.Order{{
			ID:     1,
			UserID: 10,
			Status: "pending",
			Total:  decimal.NewFromInt(200),
		}}
		getter := &mockOrdersGetter{
			listByUserIDFn: func(ctx context.Context, userID int) ([]domain.Order, error) {
				assert.Equal(t, 10, userID)
				return expected, nil
			},
		}
		svc := service.NewOrderService(&mockUnitOfWork{}, getter)

		orders, err := svc.ListByUserID(ctx, 10)
		require.NoError(t, err)
		assert.Equal(t, expected, orders)
	})

	t.Run("error", func(t *testing.T) {
		getter := &mockOrdersGetter{
			listByUserIDFn: func(ctx context.Context, userID int) ([]domain.Order, error) {
				return nil, domain.ErrNotFound
			},
		}
		svc := service.NewOrderService(&mockUnitOfWork{}, getter)

		orders, err := svc.ListByUserID(ctx, 10)
		require.ErrorIs(t, err, domain.ErrNotFound)
		assert.Empty(t, orders)
	})
}

func TestOrderService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		uow := &mockUnitOfWork{
			doFn: func(ctx context.Context, fn func(ctx context.Context, repos service.Repositories) (any, error)) (any, error) {
				// В unit-тесте не исполняем реальный callback с конкретными репозиториями —
				// проверяем, что сервис корректно обрабатывает результат UoW.
				return 77, nil
			},
		}
		svc := service.NewOrderService(uow, &mockOrdersGetter{})

		orderID, err := svc.Create(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, 77, orderID)
	})

	t.Run("uow error", func(t *testing.T) {
		uow := &mockUnitOfWork{
			doFn: func(ctx context.Context, fn func(ctx context.Context, repos service.Repositories) (any, error)) (any, error) {
				return nil, domain.ErrValidation
			},
		}
		svc := service.NewOrderService(uow, &mockOrdersGetter{})

		orderID, err := svc.Create(ctx, 1)
		require.ErrorIs(t, err, domain.ErrValidation)
		assert.Equal(t, 0, orderID)
	})

	t.Run("invalid return type", func(t *testing.T) {
		uow := &mockUnitOfWork{
			doFn: func(ctx context.Context, fn func(ctx context.Context, repos service.Repositories) (any, error)) (any, error) {
				return "not-an-int", nil
			},
		}
		svc := service.NewOrderService(uow, &mockOrdersGetter{})

		orderID, err := svc.Create(ctx, 1)
		require.ErrorIs(t, err, domain.ErrParse)
		assert.Equal(t, 0, orderID)
	})
}
