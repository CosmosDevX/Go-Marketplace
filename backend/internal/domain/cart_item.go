package domain

import "fmt"

type CartItem struct {
	ID       int
	CartID   int
	Product  Product
	Quantity int
}

func NewCartItem(cartID, productID int) (CartItem, error) {
	if cartID <= 0 {
		return CartItem{}, fmt.Errorf("cart id: %w", ErrValidation)
	}
	if productID <= 0 {
		return CartItem{}, fmt.Errorf("product id: %w", ErrValidation)
	}

	return CartItem{
		CartID: cartID,
		Product: Product{
			ID: productID,
		},
	}, nil
}
