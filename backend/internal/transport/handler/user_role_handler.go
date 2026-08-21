package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
)

type UserRoleService interface {
	Create(ctx context.Context, input service.UserRoleInput) error
	Delete(ctx context.Context, input service.UserRoleInput, currentUserID int) error
	List(ctx context.Context, username string) ([]string, error)
}

type UserRoleHandler struct {
	userRoleService UserRoleService
}

func NewUserRoleHandler(userRoleService UserRoleService) UserRoleHandler {
	return UserRoleHandler{
		userRoleService: userRoleService,
	}
}

func (h UserRoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.UserRoleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, fmt.Errorf("user role dto: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(w, fmt.Errorf("user role dto: %w", domain.ErrValidation))
		return
	}

	if err := h.userRoleService.Create(ctx, service.UserRoleInput{Username: dto.Username, Role: dto.Role}); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "role granted"})
}

func (h UserRoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	if err := h.userRoleService.Delete(ctx, service.UserRoleInput{Username: r.PathValue("username"), Role: r.PathValue("rolename")}, user.UserID); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "role deleted"})
}

func (h UserRoleHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roles, err := h.userRoleService.List(ctx, r.PathValue("username"))
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, roles)
}
