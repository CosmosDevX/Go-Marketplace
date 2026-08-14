package handler

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type UserCreater interface {
	CreateUser(ctx context.Context, userDTO dto.CreateUserDTO) (int, *domain.DomainError)
}

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

func (h UserHandler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createUserDTO dto.CreateUserDTO
	if err := utils.Deserialize(r.Body, &createUserDTO); err != nil {
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing create user dto"))
		return
	}

	if err := validator.Struct(createUserDTO); err != nil {
		utils.WriteError(w, *domain.NewValidationError(err.Error()))
		return
	}

	if err := utils.ActivateRateLimiter(ctx, w, r, constants.CreateUserRateLimitKey, &h.rateLimiter, redis_rate.PerHour(5)); err != nil {
		return
	}

	id, domainErr := h.userCreater.CreateUser(ctx, createUserDTO)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"user_id": id})
}
