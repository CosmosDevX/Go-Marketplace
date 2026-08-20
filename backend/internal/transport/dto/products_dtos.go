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
	Image       string          `json:"product_image"`
	Category    GetCategoryDTO  `json:"category"`
	SellerID    int             `json:"seller_id"`
}

type CreateProductDTO struct {
	Name         string          `json:"product_name" validate:"required,min=3,max=60"`
	Description  string          `json:"product_description" validate:"required,min=3,max=400"`
	Price        decimal.Decimal `json:"product_price" validate:"required"`
	CategorySlug string          `json:"category_slug" validate:"required,min=3,max=60"`
	SellerID     int             `json:"seller_id" validate:"required"`
}

func ToGetProductDTO(product domain.Product) GetProductDTO {
	return GetProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		Category:    ToGetCategoryDTO(product.Category),
		SellerID:    product.SellerID,
	}
}
