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

func (h ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, createProductRateLimitKey, &h.rateLimiter, redis_rate.PerHour(5)); err != nil {
		return
	}

	if err := r.ParseMultipartForm(config.MaxBodySize); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("form parse: %w", err))
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

func (h ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(config.MaxBodySize); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("form parse: %w", err))
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
