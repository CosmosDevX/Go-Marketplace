package service

import (
	"context"
	"fmt"
	"myapp/internal/domain"
)

type CreateCartItemInput struct {
	UserID    int
	ProductID int
}

type CartItemRepository interface {
	Create(ctx context.Context, cartID, productID int) (int, error)
	ListByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error)
	UpdateQuantity(ctx context.Context, cartID, cartItemID, delta int) (int, error)
	Delete(ctx context.Context, cartID, cartItemID int) error
}

type CartIDGetter interface {
	GetIDByUserID(ctx context.Context, userID int) (int, error)
}

type CartItemService struct {
	cartItemRepository CartItemRepository
	cartIDGetter       CartIDGetter
}

func NewCartItemService(cartItemRepository CartItemRepository, cartIDGetter CartIDGetter) CartItemService {
	return CartItemService{
		cartItemRepository: cartItemRepository,
		cartIDGetter:       cartIDGetter,
	}
}

func (s CartItemService) Create(ctx context.Context, input CreateCartItemInput) (int, error) {
	cartID, err := s.cartIDGetter.GetIDByUserID(ctx, input.UserID)
	if err != nil {
		return 0, err
	}

	cartItem, err := domain.NewCartItem(cartID, input.ProductID)
	if err != nil {
		return 0, err
	}

	cartItemID, err := s.cartItemRepository.Create(ctx, cartItem.CartID, cartItem.Product.ID)
	if err != nil {
		return 0, err
	}

	return cartItemID, nil
}

func (s CartItemService) List(ctx context.Context, userID int) ([]domain.CartItem, error) {
	cartID, err := s.cartIDGetter.GetIDByUserID(ctx, userID)
	if err != nil {
		return []domain.CartItem{}, err
	}

	cartItems, err := s.cartItemRepository.ListByCartID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (s CartItemService) UpdateQuantity(ctx context.Context, userID, cartItemID, delta int) (int, error) {
	if delta != 1 && delta != -1 {
		return 0, fmt.Errorf("invalid delta value: %w", domain.ErrValidation)
	}

	cartID, err := s.cartIDGetter.GetIDByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	quantity, err := s.cartItemRepository.UpdateQuantity(ctx, cartID, cartItemID, delta)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func (s CartItemService) Delete(ctx context.Context, userID, cartItemID int) error {
	cartID, err := s.cartIDGetter.GetIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.cartItemRepository.Delete(ctx, cartID, cartItemID); err != nil {
		return err
	}

	return nil
}
