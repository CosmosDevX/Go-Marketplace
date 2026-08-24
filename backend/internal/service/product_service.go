package service

import (
	"context"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"
	"slices"

	"github.com/shopspring/decimal"
)

type ProductInput struct {
	ID           int
	Name         string
	Description  string
	Price        decimal.Decimal
	CategorySlug string
	SellerID     int
	Filename     string
}

type ProductRepository interface {
	ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error)
	List(ctx context.Context, search, categorySlug, sortBy string, asc bool, page int) ([]domain.Product, error)
	Create(ctx context.Context, p domain.Product) (int, error)
	GetImageByID(ctx context.Context, productID, sellerID int, isAdmin bool) (string, error)
	Delete(ctx context.Context, productID, sellerID int, isAdmin bool) error
	Update(ctx context.Context, p domain.Product, productID int) error
}

type CategoryIDGetter interface {
	GetIDBySlug(ctx context.Context, categorySlug string) (int, error)
}

type ProductService struct {
	productRepository ProductRepository
	categoryIDGetter  CategoryIDGetter
}

func NewProductService(productRepository ProductRepository, categoryIDGetter CategoryIDGetter) ProductService {
	return ProductService{
		productRepository: productRepository,
		categoryIDGetter:  categoryIDGetter,
	}
}

func (s ProductService) Create(ctx context.Context, input ProductInput) (int, error) {
	categoryID, err := s.categoryIDGetter.GetIDBySlug(ctx, input.CategorySlug)
	if err != nil {
		return 0, err
	}

	product, err := domain.NewProduct(input.Name, input.Description, input.Filename, input.Price, categoryID, input.SellerID)
	if err != nil {
		return 0, err
	}

	productID, err := s.productRepository.Create(ctx, product)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

func (s ProductService) Update(ctx context.Context, input ProductInput) (string, error) {
	categoryID, err := s.categoryIDGetter.GetIDBySlug(ctx, input.CategorySlug)
	if err != nil {
		return "", err
	}

	product, err := domain.NewProduct(input.Name, input.Description, input.Filename, input.Price, categoryID, input.SellerID)
	if err != nil {
		return "", err
	}

	oldFilename, err := s.productRepository.GetImageByID(ctx, input.ID, input.SellerID, false)
	if err != nil {
		return "", err
	}

	if err := s.productRepository.Update(ctx, product, input.ID); err != nil {
		return "", err
	}

	return oldFilename, nil
}

func (s ProductService) List(ctx context.Context, search, categorySlug, sortBy string, asc bool, page int) ([]domain.Product, error) {
	if page <= 0 {
		return nil, fmt.Errorf("list products: %w", domain.ErrValidation)
	}

	products, err := s.productRepository.List(ctx, search, categorySlug, sortBy, asc, page)
	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (s ProductService) ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error) {
	if page <= 0 {
		return nil, fmt.Errorf("get products by sellerID: %w", domain.ErrValidation)
	}

	products, err := s.productRepository.ListBySellerID(ctx, sellerID, page)
	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (s ProductService) Delete(ctx context.Context, productID, sellerID int, roles []string) (string, error) {
	isAdmin := false
	if slices.Contains(roles, config.AdminRole) {
		isAdmin = true
	}

	filename, err := s.productRepository.GetImageByID(ctx, productID, sellerID, isAdmin)
	if err != nil {
		return "", err
	}

	if err := s.productRepository.Delete(ctx, productID, sellerID, isAdmin); err != nil {
		return "", err
	}

	return filename, nil
}
