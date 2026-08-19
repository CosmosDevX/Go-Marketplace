package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type ProductService interface {
	Create(ctx context.Context, input service.CreateProductInput) (int, error)
	List(ctx context.Context, categorySlug string, page int) ([]domain.Product, error)
}

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(productService ProductService) ProductHandler {
	return ProductHandler{
		productService: productService,
	}
}

func (h ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create product dto: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create product dto: %w", domain.ErrValidation))
		return
	}

	productID, err := h.productService.Create(ctx, service.CreateProductInput{
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
		Image:        dto.Image,
		CategorySlug: dto.CategorySlug,
	})
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"product_id": productID})
}

func (h ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categorySlug := r.URL.Query().Get("category")
	page, parseErr := strconv.Atoi(r.URL.Query().Get("page"))
	if parseErr != nil {
		utils.WriteError(w, fmt.Errorf("page: %w", domain.ErrParse))
		return
	}

	products, err := h.productService.List(ctx, categorySlug, page)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	dtos := make([]dto.GetProductDTO, len(products))
	for i := range dtos {
		dtos[i] = dto.ToGetProductDTO(products[i])
	}

	utils.WriteJSON(w, map[string]any{
		"products": dtos,
		"page":     page,
		"limit":    config.ProductPageSize,
	})
}
