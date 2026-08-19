package dto

import (
	"myapp/internal/domain"
)

type GetCartItemDTO struct {
	ID       int           `json:"cart_item_id"`
	CartID   int           `json:"cart_id"`
	Product  GetProductDTO `json:"product"`
	Quantity int           `json:"quantity"`
}

type CreateCartItemDTO struct {
	ProductID int `json:"product_id" validate:"required"`
}

func ToGetCartItemDTO(cartItem domain.CartItem) GetCartItemDTO {
	return GetCartItemDTO{
		ID:       cartItem.ID,
		CartID:   cartItem.CartID,
		Quantity: cartItem.Quantity,
		Product:  ToGetProductDTO(cartItem.Product),
	}
}
