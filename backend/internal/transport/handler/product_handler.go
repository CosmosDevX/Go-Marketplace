package handler

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type ProductService interface {
	CreateProduct(ctx context.Context, createProductDTO dto.CreateProductDTO) (int, *domain.DomainError)
	GetProductsByCategory(ctx context.Context, categorySlug string, page int) ([]dto.GetProductDTO, *domain.DomainError)
}

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(productService ProductService) ProductHandler {
	return ProductHandler{
		productService: productService,
	}
}

func (h ProductHandler) HandleProductCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createProductDTO dto.CreateProductDTO
	if err := utils.Deserialize(r.Body, &createProductDTO); err != nil {
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing create product dto"))
		return
	}

	if err := validator.Struct(createProductDTO); err != nil {
		utils.WriteError(w, *domain.NewValidationError(err.Error()))
		return
	}

	productID, domainErr := h.productService.CreateProduct(ctx, createProductDTO)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"product_id": productID})
}

func (h ProductHandler) HandleGetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categorySlug := r.URL.Query().Get("category")
	page, parseErr := strconv.Atoi(r.URL.Query().Get("page"))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewDomainError(constants.ParseError, "error during parsing page number"))
		return
	}

	getProductDTOs, domainErr := h.productService.GetProductsByCategory(ctx, categorySlug, page)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]any{
		"products": getProductDTOs,
		"page":     page,
		"limit":    16,
	})
}
