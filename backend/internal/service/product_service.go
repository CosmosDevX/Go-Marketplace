package service

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
)

type ProductRepository interface {
	ListByCategorySlug(ctx context.Context, categorySlug string, page int) ([]domain.Product, error)
	ListAll(ctx context.Context, page int) ([]domain.Product, error)
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

func (s ProductService) Create(ctx context.Context, dto dto.CreateProductDTO) (int, error) {
	categoryID, err := s.categoryIDGetter.GetIDBySlug(ctx, dto.CategorySlug)
	if err != nil {
		return 0, err
	}

	product, err := domain.NewProduct(dto.Name, dto.Description, dto.Image, dto.Price, dto.Quantity, categoryID)
	if err != nil {
		return 0, err
	}

	productID, err := s.productRepository.Create(ctx, product)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

func (s ProductService) List(ctx context.Context, categorySlug string, page int) ([]dto.GetProductDTO, error) {
	if page <= 0 {
		return nil, fmt.Errorf("get products by category slug: %w", domain.ErrValidation)
	}

	var products []domain.Product
	var err error
	if categorySlug == "" {
		products, err = s.productRepository.ListAll(ctx, page)
	} else {
		products, err = s.productRepository.ListByCategorySlug(ctx, categorySlug, page)
	}

	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return []dto.GetProductDTO{}, nil
	}

	dtos := make([]dto.GetProductDTO, len(products))
	for i := range dtos {
		dtos[i] = dto.ToGetProductDTO(products[i])
	}

	return dtos, nil
}
