package handler

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
)

type CategoryService interface {
	Create(ctx context.Context, createCategoryDTO dto.CreateCategoryDTO) (int, error)
	ListAll(ctx context.Context) ([]dto.GetCategoryDTO, error)
}

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) CategoryHandler {
	return CategoryHandler{
		categoryService: categoryService,
	}
}

func (h CategoryHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createCategoryDTO dto.CreateCategoryDTO
	if err := utils.Deserialize(r.Body, &createCategoryDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("create category dto deserialize: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(createCategoryDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("create category dto validate: %w", domain.ErrValidation))
		return
	}

	categoryID, err := h.categoryService.Create(ctx, createCategoryDTO)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"category_id": categoryID})
}

func (h CategoryHandler) ListAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	getCategoryDTOs, err := h.categoryService.ListAll(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, getCategoryDTOs)
}
