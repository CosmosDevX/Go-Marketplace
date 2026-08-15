package service

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/internal/transport/dto"
)

type ProductRepository interface {
	ListByCategorySlug(ctx context.Context, categorySlug string, page int) ([]domain.Product, error)
	ListAll(ctx context.Context, page int) ([]domain.Product, error)
	Create(ctx context.Context, params repository.CreateProductParams) (int, error)
}

type ProductService struct {
	unitOfWork        UnitOfWork
	productRepository ProductRepository
}

func NewProductService(unitOfWork UnitOfWork, productRepository ProductRepository) ProductService {
	return ProductService{
		unitOfWork:        unitOfWork,
		productRepository: productRepository,
	}
}

func (s ProductService) Create(ctx context.Context, createProductDTO dto.CreateProductDTO) (int, error) {
	value, err := s.unitOfWork.Do(ctx, func(ctx context.Context, repos Repositories) (any, error) {
		categoryID, err := repos.CategoryRepository.GetIDBySlug(ctx, createProductDTO.CategorySlug)
		if err != nil {
			return 0, err
		}

		productID, err := repos.ProductRepository.Create(ctx, repository.CreateProductParams{
			Name:        createProductDTO.Name,
			Description: createProductDTO.Description,
			Price:       createProductDTO.Price,
			Image:       createProductDTO.Image,
			CategoryID:  categoryID,
		})
		if err != nil {
			return 0, err
		}

		return productID, nil
	})

	if err != nil {
		return 0, err
	}

	productID, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("product id parse: %w", domain.ErrParse)
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
