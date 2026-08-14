package service

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/internal/transport/dto"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]domain.Category, *domain.DomainError)
	CreateCategory(ctx context.Context, params repository.CreateCategoryParams) (int, *domain.DomainError)
}

type CategoryService struct {
	categoryRepository CategoryRepository
}

func NewCategoryService(categoryRepository CategoryRepository) CategoryService {
	return CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s CategoryService) CreateCategory(ctx context.Context, createCategoryDTO dto.CreateCategoryDTO) (int, *domain.DomainError) {
	categoryID, domainErr := s.categoryRepository.CreateCategory(ctx, repository.CreateCategoryParams{Name: createCategoryDTO.Name, Slug: createCategoryDTO.Slug})
	if domainErr != nil {
		return 0, domainErr
	}

	return categoryID, nil
}

func (s CategoryService) GetAllCategories(ctx context.Context) ([]dto.GetCategoryDTO, *domain.DomainError) {
	categories, domainErr := s.categoryRepository.GetAllCategories(ctx)
	if domainErr != nil {
		return nil, domainErr
	}
	if len(categories) == 0 {
		return nil, nil
	}

	dtos := make([]dto.GetCategoryDTO, len(categories))
	for i := range dtos {
		dtos[i] = dto.ToGetCategoryDTO(categories[i])
	}

	return dtos, nil
}
