package service

import (
	"context"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type ordersGetter interface {
	ListByUserID(ctx context.Context, userID int) ([]domain.Order, error)
}

type OrderService struct {
	unitOfWork   UnitOfWork
	ordersGetter ordersGetter
}

func NewOrderService(unitOfWork UnitOfWork, ordersGetter ordersGetter) OrderService {
	return OrderService{
		unitOfWork:   unitOfWork,
		ordersGetter: ordersGetter,
	}
}

func (s OrderService) Create(ctx context.Context, userID int) (int, error) {
	value, err := s.unitOfWork.Do(ctx, func(ctx context.Context, repos Repositories) (any, error) {
		var total decimal.Decimal
		cartItemsTotal := map[int]decimal.Decimal{}

		cartID, err := repos.CartRepository.GetIDByUserID(ctx, userID)
		if err != nil {
			return 0, err
		}

		cartItems, err := repos.CartItemRepository.ListByCartID(ctx, cartID)
		if err != nil {
			return 0, err
		}

		if len(cartItems) == 0 {
			return 0, fmt.Errorf("cart is empty: %w", domain.ErrValidation)
		}

		for i := range cartItems {
			quantity := decimal.NewFromInt(int64(cartItems[i].Quantity))
			multiple := cartItems[i].Product.Price.Mul(quantity)
			total = total.Add(multiple)
			cartItemsTotal[cartItems[i].ID] = multiple
		}

		orderID, err := repos.OrderRepository.Create(ctx, userID, config.PendingOrderStatus, total)
		if err != nil {
			return 0, err
		}

		var domainModels []domain.OrderItem
		for i := range cartItems {
			domainModel, err := domain.NewOrderItem(orderID, cartItems[i].Product.ID, cartItems[i].Product.SellerID, cartItems[i].Quantity,
				cartItemsTotal[cartItems[i].ID], cartItems[i].Product.Name, cartItems[i].Product.Price)
			if err != nil {
				return nil, err
			}

			domainModels = append(domainModels, domainModel)
		}

		if err := repos.OrderItemRepository.CreateMany(ctx, domainModels); err != nil {
			return 0, err
		}

		if err := repos.CartItemRepository.DeleteAllByCartID(ctx, cartID); err != nil {
			return 0, err
		}

		return orderID, nil
	})

	if err != nil {
		return 0, err
	}

	orderID, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("orderID parse: %w", domain.ErrParse)
	}

	return orderID, nil
}

func (s OrderService) ListByUserID(ctx context.Context, userID int) ([]domain.Order, error) {
	orders, err := s.ordersGetter.ListByUserID(ctx, userID)
	if err != nil {
		return []domain.Order{}, err
	}

	return orders, nil
}
