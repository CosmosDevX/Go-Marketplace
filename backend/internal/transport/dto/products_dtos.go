package dto

import (
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type GetProductDTO struct {
	ID          int             `json:"product_id"`
	Name        string          `json:"product_name"`
	Description string          `json:"product_description"`
	Price       decimal.Decimal `json:"product_price"`
	Quantity    int             `json:"product_quantity"`
	Image       string          `json:"product_image"`
	Category    GetCategoryDTO  `json:"category"`
}

type CreateProductDTO struct {
	Name         string          `json:"product_name" validate:"required,min=3,max=60"`
	Description  string          `json:"product_description" validate:"required,min=3,max=400"`
	Price        decimal.Decimal `json:"product_price" validate:"required"`
	Quantity     int             `json:"product_quantity" validate:"gte=0"`
	Image        string          `json:"product_image" validate:"required"`
	CategorySlug string          `json:"category_slug" validate:"required,min=5,max=60"`
}

func ToGetProductDTO(product domain.Product) GetProductDTO {
	return GetProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Image:       product.Image,
		Category:    ToGetCategoryDTO(product.Category),
	}
}
