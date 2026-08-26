package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID         int
	UserID     int
	Status     string
	OrderItems []OrderItem
	Total      decimal.Decimal
}

func NewOrder(userID int, status string, total decimal.Decimal) (Order, error) {
	if userID <= 0 {
		return Order{}, fmt.Errorf("userID: %w", ErrValidation)
	}
	if status == "" {
		return Order{}, fmt.Errorf("status: %w", ErrValidation)
	}
	if total.LessThanOrEqual(decimal.Zero) {
		return Order{}, fmt.Errorf("total: %w", ErrValidation)
	}

	return Order{
		UserID: userID,
		Status: status,
		Total:  total,
	}, nil
}
