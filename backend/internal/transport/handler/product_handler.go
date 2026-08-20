package handler

import (
	"context"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type ProductService interface {
	Create(ctx context.Context, input service.CreateProductInput) (int, error)
	ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error)
	List(ctx context.Context, categorySlug string, page int) ([]domain.Product, error)
	Delete(ctx context.Context, productID, sellerID int) error
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

	if err := r.ParseMultipartForm(20 * 1024 * 1024); err != nil {
		utils.WriteError(w, fmt.Errorf("form parse: %w", err))
		return
	}

	file, header, err := r.FormFile("product_image")
	if err != nil {
		utils.WriteError(w, fmt.Errorf("get file from form: %w", err))
		return
	}

	var contentType string
	if header != nil {
		contentType = header.Header.Get("Content-Type")
	}
	if file != nil && !strings.HasPrefix(contentType, "image/") {
		utils.WriteError(w, fmt.Errorf("invalid file type: %w", domain.ErrValidation))
		return
	}
	if file != nil {
		defer file.Close()
	}

	productPrice, err := decimal.NewFromString(r.FormValue("product_price"))
	if err != nil {
		utils.WriteError(w, fmt.Errorf("string to decimal product price: %w", err))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	dto := dto.CreateProductDTO{
		Name:         r.FormValue("product_name"),
		Description:  r.FormValue("product_description"),
		Price:        productPrice,
		CategorySlug: r.FormValue("category_slug"),
		SellerID:     user.UserID,
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(w, fmt.Errorf("create product dto: %w", domain.ErrValidation))
		return
	}

	productID, err := h.productService.Create(ctx, service.CreateProductInput{
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
		CategorySlug: dto.CategorySlug,
		SellerID:     dto.SellerID,
		File:         file,
		Header:       header,
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
		utils.WriteError(w, fmt.Errorf("page parse: %w", domain.ErrParse))
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

func (h ProductHandler) ListBySeller(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, parseErr := strconv.Atoi(r.URL.Query().Get("page"))
	if parseErr != nil {
		utils.WriteError(w, fmt.Errorf("page parse: %w", domain.ErrParse))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	products, err := h.productService.ListBySellerID(ctx, user.UserID, page)
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

func (h ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productID, err := strconv.Atoi(r.PathValue("product_id"))
	if err != nil {
		utils.WriteError(w, fmt.Errorf("product id parse: %w", domain.ErrParse))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	if err := h.productService.Delete(ctx, productID, user.UserID); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "delete successful"})
}
