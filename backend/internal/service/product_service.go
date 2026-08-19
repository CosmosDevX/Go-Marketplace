package service

import (
	"context"
	"fmt"
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type CreateProductInput struct {
	Name         string
	Description  string
	Price        decimal.Decimal
	Quantity     int
	Image        string
	CategorySlug string
}

type ProductRepository interface {
	ListByCategorySlug(ctx context.Context, categorySlug string, page int) ([]domain.Product, error)
	List(ctx context.Context, page int) ([]domain.Product, error)
	Create(ctx context.Context, p domain.Product) (int, error)
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

func (s ProductService) Create(ctx context.Context, input CreateProductInput) (int, error) {
	categoryID, err := s.categoryIDGetter.GetIDBySlug(ctx, input.CategorySlug)
	if err != nil {
		return 0, err
	}

	product, err := domain.NewProduct(input.Name, input.Description, input.Image, input.Price, categoryID)
	if err != nil {
		return 0, err
	}

	productID, err := s.productRepository.Create(ctx, product)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

func (s ProductService) List(ctx context.Context, categorySlug string, page int) ([]domain.Product, error) {
	if page <= 0 {
		return nil, fmt.Errorf("get products by category slug: %w", domain.ErrValidation)
	}

	var products []domain.Product
	var err error
	if categorySlug == "" {
		products, err = s.productRepository.List(ctx, page)
	} else {
		products, err = s.productRepository.ListByCategorySlug(ctx, categorySlug, page)
	}

	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}
