package dto

import (
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type GetOrderItemDTO struct {
	ID       int             `json:"order_item_id"`
	Product  GetProductDTO   `json:"product"`
	Quantity int             `json:"order_item_quantity"`
	Total    decimal.Decimal `json:"order_item_total"`
}

type GetOrderDTO struct {
	ID            int               `json:"order_id"`
	Status        string            `json:"order_status"`
	UserID        int               `json:"user_id"`
	Total         decimal.Decimal   `json:"order_total"`
	OrderItemDTOs []GetOrderItemDTO `json:"order_items"`
}

func ToGetOrderDTO(order domain.Order) GetOrderDTO {
	var orderItemDTOs []GetOrderItemDTO
	for i := range order.OrderItems {
		orderItemDTOs = append(orderItemDTOs, GetOrderItemDTO{
			ID:       order.OrderItems[i].ID,
			Product:  ToGetProductDTO(order.OrderItems[i].Product),
			Quantity: order.OrderItems[i].Quantity,
			Total:    order.OrderItems[i].Total,
		})
	}

	return GetOrderDTO{
		ID:            order.ID,
		Status:        order.Status,
		UserID:        order.UserID,
		Total:         order.Total,
		OrderItemDTOs: orderItemDTOs,
	}
}
