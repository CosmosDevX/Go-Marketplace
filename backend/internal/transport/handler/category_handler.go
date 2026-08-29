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

// Create godoc
//
//	@Summary		Create category
//	@Description	Create a new product category. Requires admin role.
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateCategoryDTO	true	"Category data"
//	@Success		200		{object}	CategoryIDResponse		"Created category ID"
//	@Failure		400		{object}	ErrorResponse			"Parse or validation error"
//	@Failure		401		{object}	ErrorResponse			"Unauthorized"
//	@Failure		403		{object}	ErrorResponse			"Forbidden (not admin)"
//	@Failure		409		{object}	ErrorResponse			"Category name or slug already exists"
//	@Failure		500		{object}	ErrorResponse			"Internal server error"
//	@Router			/categories [post]
func (h CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("create category dto: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("create category dto: %w", domain.ErrValidation))
		return
	}

	categoryID, err := h.categoryService.Create(ctx, service.CreateCategoryInput{
		Name: dto.Name,
		Slug: dto.Slug,
	})
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"category_id": categoryID})
}

// List godoc
//
//	@Summary		List categories
//	@Description	Get all product categories.
//	@Tags			Categories
//	@Produce		json
//	@Success		200	{array}		dto.GetCategoryDTO	"List of categories"
//	@Failure		500	{object}	ErrorResponse		"Internal server error"
//	@Router			/categories [get]
func (h CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.categoryService.List(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	dtos := make([]dto.GetCategoryDTO, len(categories))
	for i := range dtos {
		dtos[i] = dto.ToGetCategoryDTO(categories[i])
	}

	utils.WriteJSON(w, dtos)
}
