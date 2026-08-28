package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type UserCreator interface {
	Create(ctx context.Context, input service.CreateUserInput) (int, error)
}

const createUserRateLimitKey = "createUser"

type UserHandler struct {
	userCreator UserCreator
	rateLimiter redis_rate.Limiter
}

func NewUserHandler(userCreator UserCreator, rateLimiter redis_rate.Limiter) UserHandler {
	return UserHandler{
		userCreator: userCreator,
		rateLimiter: rateLimiter,
	}
}

func (h UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, createUserRateLimitKey, &h.rateLimiter, redis_rate.PerHour(5)); err != nil {
		return
	}

	var dto dto.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create user dto: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create user dto: %w", domain.ErrValidation))
		return
	}

	id, err := h.userCreator.Create(ctx, service.CreateUserInput{
		Username: dto.Username,
		Password: dto.Password,
		Email:    dto.Email,
	})
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"user_id": id})
}
