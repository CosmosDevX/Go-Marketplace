package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID       int
	OrderID  int
	Product  Product
	Quantity int
	Total    decimal.Decimal
}

func NewOrderItem(orderID, productID, sellerID, quantity int, total decimal.Decimal, productName string, productPrice decimal.Decimal) (OrderItem, error) {
	if orderID <= 0 {
		return OrderItem{}, fmt.Errorf("orderID: %w", ErrValidation)
	}
	if productID <= 0 {
		return OrderItem{}, fmt.Errorf("cartItemID: %w", ErrValidation)
	}
	if sellerID <= 0 {
		return OrderItem{}, fmt.Errorf("sellerID: %w", ErrValidation)
	}
	if quantity <= 0 {
		return OrderItem{}, fmt.Errorf("quantity: %w", ErrValidation)
	}
	if productPrice.LessThanOrEqual(decimal.Zero) {
		return OrderItem{}, fmt.Errorf("product price: %w", ErrValidation)
	}
	if total.LessThanOrEqual(decimal.Zero) {
		return OrderItem{}, fmt.Errorf("total: %w", ErrValidation)
	}
	if productName == "" {
		return OrderItem{}, fmt.Errorf("product name: %w", ErrValidation)
	}

	return OrderItem{
		OrderID: orderID,
		Product: Product{
			ID:       productID,
			Name:     productName,
			Price:    productPrice,
			SellerID: sellerID,
		},
		Quantity: quantity,
		Total:    total,
	}, nil
}
