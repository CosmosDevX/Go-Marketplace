package service

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/internal/transport/dto"
)

type ProductRepository interface {
	GetProductsByCategory(ctx context.Context, categorySlug string, page int) ([]domain.Product, *domain.DomainError)
	CreateProduct(ctx context.Context, params repository.CreateProductParams) (int, *domain.DomainError)
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

func (s ProductService) CreateProduct(ctx context.Context, createProductDTO dto.CreateProductDTO) (int, *domain.DomainError) {
	value, domainErr := s.unitOfWork.Do(ctx, func(ctx context.Context, repos Repositories) (any, *domain.DomainError) {
		categoryID, domainErr := repos.CategoryRepository.GetCategoryIDBySlug(ctx, createProductDTO.CategorySlug)
		if domainErr != nil {
			return 0, domainErr
		}

		productID, domainErr := repos.ProductRepository.CreateProduct(ctx, repository.CreateProductParams{
			Name:        createProductDTO.Name,
			Description: createProductDTO.Description,
			Price:       createProductDTO.Price,
			Image:       createProductDTO.Image,
			CategoryID:  categoryID,
		})
		if domainErr != nil {
			return 0, domainErr
		}

		return productID, nil
	})

	if domainErr != nil {
		return 0, domainErr
	}
	productID := value.(int)

	return productID, nil
}

func (s ProductService) GetProductsByCategory(ctx context.Context, categorySlug string, page int) ([]dto.GetProductDTO, *domain.DomainError) {
	products, domainErr := s.productRepository.GetProductsByCategory(ctx, categorySlug, page)
	if domainErr != nil {
		return nil, domainErr
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
