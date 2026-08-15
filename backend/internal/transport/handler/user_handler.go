package handler

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type UserCreater interface {
	Create(ctx context.Context, userDTO dto.CreateUserDTO) (int, error)
}

const createUserRateLimitKey = "create_user"

type UserHandler struct {
	userCreater UserCreater
	rateLimiter redis_rate.Limiter
}

func NewUserHandler(userCreater UserCreater, rateLimiter redis_rate.Limiter) UserHandler {
	return UserHandler{
		userCreater: userCreater,
		rateLimiter: rateLimiter,
	}
}

func (h UserHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createUserDTO dto.CreateUserDTO
	if err := utils.Deserialize(r.Body, &createUserDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("create user dto deserialize: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(createUserDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("create user dto validate: %w", domain.ErrValidation))
		return
	}

	if err := utils.ActivateRateLimiter(ctx, w, r, createUserRateLimitKey, &h.rateLimiter, redis_rate.PerHour(5)); err != nil {
		return
	}

	id, err := h.userCreater.Create(ctx, createUserDTO)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"user_id": id})
}
