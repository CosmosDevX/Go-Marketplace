package handler

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"myapp/internal/config"
	"myapp/internal/domain"
	"myapp/internal/logger"
	"myapp/internal/service"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis_rate/v10"
	"github.com/shopspring/decimal"
)

type ProductService interface {
	Create(ctx context.Context, input service.ProductInput) (int, error)
	ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error)
	List(ctx context.Context, search, categorySlug, sortBy string, asc bool, page int) ([]domain.Product, error)
	Delete(ctx context.Context, productID, sellerID int, roles []string) (string, error)
	Update(ctx context.Context, input service.ProductInput) (string, error)
}

type FileManager interface {
	SaveFile(file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteFile(filename string) error
}

const createProductRateLimitKey = "createProduct"

type ProductHandler struct {
	productService ProductService
	fileManager    FileManager
	rateLimiter    redis_rate.Limiter
}

func NewProductHandler(productService ProductService, fileManager FileManager, rateLimiter redis_rate.Limiter) ProductHandler {
	return ProductHandler{
		productService: productService,
		fileManager:    fileManager,
		rateLimiter:    rateLimiter,
	}
}

// Create godoc
//
//	@Summary		Create product
//	@Description	Create a new product with image upload. Requires seller or admin role. Rate limit: 5 req/hour. Content-Type: multipart/form-data.
//	@Tags			SellerProducts
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product_name		formData	string	true	"Product name"
//	@Param			product_description	formData	string	true	"Product description"
//	@Param			product_price		formData	string	true	"Product price (decimal)"
//	@Param			category_slug		formData	string	true	"Category slug"
//	@Param			product_image		formData	file	true	"Product image (image/*)"
//	@Success		200					{object}	ProductIDResponse	"Created product ID"
//	@Failure		400					{object}	ErrorResponse		"Parse, validation or missing file error"
//	@Failure		401					{object}	ErrorResponse		"Unauthorized"
//	@Failure		403					{object}	ErrorResponse		"Forbidden"
//	@Failure		429					{object}	ErrorResponse		"Too many requests"
//	@Failure		500					{object}	ErrorResponse		"Internal server error"
//	@Router			/seller/products [post]
func (h ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, createProductRateLimitKey, &h.rateLimiter, redis_rate.PerHour(5)); err != nil {
		return
	}

	if err := r.ParseMultipartForm(config.MaxBodySize); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("form parse: %w", domain.ErrParse))
		return
	}

	parsedFile, err := h.parseImage(r, "product_image")
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	dto, err := h.makeProduct(ctx, r)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	filename, err := h.fileManager.SaveFile(parsedFile.File, parsedFile.Header)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}
	parsedFile.File.Close()

	productID, err := h.productService.Create(ctx, service.ProductInput{
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
		CategorySlug: dto.CategorySlug,
		SellerID:     dto.SellerID,
		Filename:     filename,
	})
	if err != nil {
		if delErr := h.fileManager.DeleteFile(filename); delErr != nil {
			logger.FromContext(ctx).Error("failed to delete product image after create error", "error", delErr)
		}
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"product_id": productID})
}

// Update godoc
//
//	@Summary		Update product
//	@Description	Update existing product. Image is optional. Requires seller or admin role. Content-Type: multipart/form-data.
//	@Tags			SellerProducts
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product_id			path		int		true	"Product ID"
//	@Param			product_name		formData	string	true	"Product name"
//	@Param			product_description	formData	string	true	"Product description"
//	@Param			product_price		formData	string	true	"Product price (decimal)"
//	@Param			category_slug		formData	string	true	"Category slug"
//	@Param			product_image		formData	file	false	"New product image (optional)"
//	@Success		200					{object}	ProductIDResponse	"Updated product ID"
//	@Failure		400					{object}	ErrorResponse		"Parse or validation error"
//	@Failure		401					{object}	ErrorResponse		"Unauthorized"
//	@Failure		403					{object}	ErrorResponse		"Forbidden"
//	@Failure		404					{object}	ErrorResponse		"Product not found"
//	@Failure		500					{object}	ErrorResponse		"Internal server error"
//	@Router			/seller/products/{product_id} [put]
func (h ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(config.MaxBodySize); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("form parse: %w", domain.ErrParse))
		return
	}

	parsedFile, err := h.parseImage(r, "product_image")
	if err != nil {
		if !errors.Is(err, domain.ErrMissingFile) {
			utils.WriteError(ctx, w, err)
			return
		}
	}

	dto, err := h.makeProduct(ctx, r)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	var filename string
	if parsedFile.File != nil && parsedFile.Header != nil {
		filename, err = h.fileManager.SaveFile(parsedFile.File, parsedFile.Header)
		if err != nil {
			utils.WriteError(ctx, w, err)
			return
		}
	}
	parsedFile.File.Close()

	productID, err := strconv.Atoi(r.PathValue("product_id"))
	if err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("productID parse: %w", domain.ErrParse))
		return
	}

	oldFilename, err := h.productService.Update(ctx, service.ProductInput{
		ID:           productID,
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
		CategorySlug: dto.CategorySlug,
		SellerID:     dto.SellerID,
		Filename:     filename,
	})
	if err != nil {
		if delErr := h.fileManager.DeleteFile(filename); delErr != nil {
			logger.FromContext(ctx).Error("failed to delete product image after update error", "error", delErr)
		}
		utils.WriteError(ctx, w, err)
		return
	}

	if filename != "" {
		if err := h.fileManager.DeleteFile(oldFilename); err != nil {
			utils.WriteError(ctx, w, err)
			return
		}
	}

	utils.WriteJSON(w, map[string]int{"product_id": productID})
}

// List godoc
//
//	@Summary		List products
//	@Description	Public product listing with search, category filter, sorting and pagination.
//	@Tags			Products
//	@Produce		json
//	@Param			page		query		int		true	"Page number (starting from 1)"
//	@Param			search		query		string	false	"Search by product name"
//	@Param			category	query		string	false	"Filter by category slug"
//	@Param			sortBy		query		string	false	"Sort field (e.g. price, name)"
//	@Param			asc			query		bool	false	"Ascending order (default false)"
//	@Success		200			{object}	ProductListResponse	"Paginated list of products"
//	@Failure		400			{object}	ErrorResponse			"Invalid page parameter"
//	@Failure		500			{object}	ErrorResponse			"Internal server error"
//	@Router			/products [get]
func (h ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categorySlug := r.URL.Query().Get("category")
	sortBy := r.URL.Query().Get("sortBy")
	search := r.URL.Query().Get("search")
	asc, err := strconv.ParseBool(r.URL.Query().Get("asc"))
	if err != nil {
		asc = false
	}

	page, parseErr := strconv.Atoi(r.URL.Query().Get("page"))
	if parseErr != nil {
		utils.WriteError(ctx, w, fmt.Errorf("page parse: %w", domain.ErrParse))
		return
	}

	products, err := h.productService.List(ctx, search, categorySlug, sortBy, asc, page)
	if err != nil {
		utils.WriteError(ctx, w, err)
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

// ListBySeller godoc
//
//	@Summary		List seller products
//	@Description	List products belonging to the authenticated seller. Requires seller or admin role.
//	@Tags			SellerProducts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int	true	"Page number"
//	@Success		200		{object}	ProductListResponse	"Paginated list of seller products"
//	@Failure		400		{object}	ErrorResponse			"Invalid page"
//	@Failure		401		{object}	ErrorResponse			"Unauthorized"
//	@Failure		403		{object}	ErrorResponse			"Forbidden"
//	@Failure		500		{object}	ErrorResponse			"Internal server error"
//	@Router			/seller/products [get]
func (h ProductHandler) ListBySeller(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, parseErr := strconv.Atoi(r.URL.Query().Get("page"))
	if parseErr != nil {
		utils.WriteError(ctx, w, fmt.Errorf("page parse: %w", domain.ErrParse))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	products, err := h.productService.ListBySellerID(ctx, user.UserID, page)
	if err != nil {
		utils.WriteError(ctx, w, err)
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

// Delete godoc
//
//	@Summary		Delete product
//	@Description	Delete product by ID. Seller can delete own products, admin can delete any. Requires seller or admin role.
//	@Tags			SellerProducts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	MessageResponse	"Deletion confirmation"
//	@Failure		400			{object}	ErrorResponse		"Invalid product_id"
//	@Failure		401			{object}	ErrorResponse		"Unauthorized"
//	@Failure		403			{object}	ErrorResponse		"Forbidden"
//	@Failure		404			{object}	ErrorResponse		"Product not found"
//	@Failure		500			{object}	ErrorResponse		"Internal server error"
//	@Router			/seller/products/{product_id} [delete]
func (h ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productID, err := strconv.Atoi(r.PathValue("product_id"))
	if err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("product id parse: %w", domain.ErrParse))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	filename, err := h.productService.Delete(ctx, productID, user.UserID, user.Roles)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	if err := h.fileManager.DeleteFile(filename); err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "delete successful"})
}

type ParsedFile struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func (h ProductHandler) parseImage(r *http.Request, formFilename string) (ParsedFile, error) {
	file, header, err := r.FormFile(formFilename)
	if file == nil || header == nil {
		return ParsedFile{}, fmt.Errorf("file is missing: %w", domain.ErrMissingFile)
	}
	if err != nil {
		return ParsedFile{}, fmt.Errorf("get file from form: %w", err)
	}

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return ParsedFile{}, fmt.Errorf("invalid file type: %w", domain.ErrValidation)
	}

	return ParsedFile{File: file, Header: header}, nil
}

func (h ProductHandler) makeProduct(ctx context.Context, r *http.Request) (dto.CreateProductDTO, error) {
	productPrice, err := decimal.NewFromString(r.FormValue("product_price"))
	if err != nil {
		return dto.CreateProductDTO{}, fmt.Errorf("string to decimal product price: %w", err)
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		return dto.CreateProductDTO{}, err
	}

	productDTO := dto.CreateProductDTO{
		Name:         r.FormValue("product_name"),
		Description:  r.FormValue("product_description"),
		Price:        productPrice,
		CategorySlug: r.FormValue("category_slug"),
		SellerID:     user.UserID,
	}

	if err := validator.Struct(productDTO); err != nil {
		return dto.CreateProductDTO{}, fmt.Errorf("create product dto: %w", domain.ErrValidation)
	}

	return productDTO, nil
}
