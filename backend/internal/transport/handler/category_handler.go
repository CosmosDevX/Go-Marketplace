package handler

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, createCategoryDTO dto.CreateCategoryDTO) (int, *domain.DomainError)
	GetAllCategories(ctx context.Context) ([]dto.GetCategoryDTO, *domain.DomainError)
}

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) CategoryHandler {
	return CategoryHandler{
		categoryService: categoryService,
	}
}

func (h CategoryHandler) HandleCategoryCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createCategoryDTO dto.CreateCategoryDTO
	if err := utils.Deserialize(r.Body, &createCategoryDTO); err != nil {
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing create category dto"))
		return
	}

	if err := validator.Struct(createCategoryDTO); err != nil {
		utils.WriteError(w, *domain.NewValidationError(err.Error()))
		return
	}

	categoryID, domainErr := h.categoryService.CreateCategory(ctx, createCategoryDTO)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"category_id": categoryID})
}

func (h CategoryHandler) HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	getCategoryDTOs, domainErr := h.categoryService.GetAllCategories(ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, getCategoryDTOs)
}
