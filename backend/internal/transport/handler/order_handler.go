package handler

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type orderService interface {
	Create(ctx context.Context, userID int) (int, error)
	ListByUserID(ctx context.Context, userID int) ([]domain.Order, error)
}

const createOrderRateLimitKey = "createOrder"

type OrderHandler struct {
	orderService orderService
	rateLimiter  redis_rate.Limiter
}

func NewOrderHandler(orderService orderService, rateLimiter redis_rate.Limiter) OrderHandler {
	return OrderHandler{
		orderService: orderService,
		rateLimiter:  rateLimiter,
	}
}

func (h OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, createOrderRateLimitKey, &h.rateLimiter, redis_rate.PerMinute(2)); err != nil {
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	orderID, err := h.orderService.Create(ctx, user.UserID)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"order_id": orderID})
}

func (h OrderHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	orders, err := h.orderService.ListByUserID(ctx, user.UserID)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	var dtos []dto.GetOrderDTO
	for i := range orders {
		dtos = append(dtos, dto.ToGetOrderDTO(orders[i]))
	}

	utils.WriteJSON(w, dtos)
}
