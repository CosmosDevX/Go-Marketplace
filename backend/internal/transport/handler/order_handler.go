package handler

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/utils"
	"net/http"
)

type orderService interface {
	Create(ctx context.Context, userID int) (int, error)
	ListByUserID(ctx context.Context, userID int) ([]domain.Order, error)
}

type OrderHandler struct {
	orderService orderService
}

func NewOrderHandler(orderService orderService) OrderHandler {
	return OrderHandler{
		orderService: orderService,
	}
}

func (h OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	orderID, err := h.orderService.Create(ctx, user.UserID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"order_id": orderID})
}

func (h OrderHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	orders, err := h.orderService.ListByUserID(ctx, user.UserID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	var dtos []dto.GetOrderDTO
	for i := range orders {
		dtos = append(dtos, dto.ToGetOrderDTO(orders[i]))
	}

	utils.WriteJSON(w, dtos)
}
