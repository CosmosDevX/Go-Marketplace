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
)

type CategoryService interface {
	Create(ctx context.Context, input service.CreateCategoryInput) (int, error)
	List(ctx context.Context) ([]domain.Category, error)
}

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) CategoryHandler {
	return CategoryHandler{
		categoryService: categoryService,
	}
}

func (h CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create category dto: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create category dto: %w", domain.ErrValidation))
		return
	}

	categoryID, err := h.categoryService.Create(ctx, service.CreateCategoryInput{
		Name: dto.Name,
		Slug: dto.Slug,
	})
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"category_id": categoryID})
}

func (h CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.categoryService.List(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	dtos := make([]dto.GetCategoryDTO, len(categories))
	for i := range dtos {
		dtos[i] = dto.ToGetCategoryDTO(categories[i])
	}

	utils.WriteJSON(w, dtos)
}
