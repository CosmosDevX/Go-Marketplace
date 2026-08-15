package handler

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type ProductService interface {
	Create(ctx context.Context, createProductDTO dto.CreateProductDTO) (int, error)
	List(ctx context.Context, categorySlug string, page int) ([]dto.GetProductDTO, error)
}

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(productService ProductService) ProductHandler {
	return ProductHandler{
		productService: productService,
	}
}

func (h ProductHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createProductDTO dto.CreateProductDTO
	if err := utils.Deserialize(r.Body, &createProductDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("create product dto deserialize: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(createProductDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("create product dto validate: %w", domain.ErrValidation))
		return
	}

	productID, err := h.productService.Create(ctx, createProductDTO)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"product_id": productID})
}

func (h ProductHandler) ListAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categorySlug := r.URL.Query().Get("category")
	page, parseErr := strconv.Atoi(r.URL.Query().Get("page"))
	if parseErr != nil {
		utils.WriteError(w, fmt.Errorf("page parse: %w", domain.ErrParse))
		return
	}

	getProductDTOs, err := h.productService.List(ctx, categorySlug, page)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]any{
		"products": getProductDTOs,
		"page":     page,
		"limit":    len(getProductDTOs),
	})
}
