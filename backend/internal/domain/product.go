package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       decimal.Decimal
	Image       string
	Category    Category
	SellerID    int
}

func NewProduct(name, description, image string, price decimal.Decimal, categoryID, sellerID int) (Product, error) {
	if name == "" {
		return Product{}, fmt.Errorf("name: %w", ErrValidation)
	}
	if description == "" {
		return Product{}, fmt.Errorf("description: %w", ErrValidation)
	}
	if price.LessThanOrEqual(decimal.Zero) {
		return Product{}, fmt.Errorf("price must be positive: %w", ErrValidation)
	}
	if categoryID <= 0 {
		return Product{}, fmt.Errorf("category id: %w", ErrValidation)
	}
	if sellerID <= 0 {
		return Product{}, fmt.Errorf("seller id: %w", ErrValidation)
	}

	return Product{
		Name:        name,
		Description: description,
		Image:       image,
		Price:       price,
		Category: Category{
			ID: categoryID,
		},
		SellerID: sellerID,
	}, nil
}
