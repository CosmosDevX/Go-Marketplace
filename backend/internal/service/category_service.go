package service

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/internal/transport/dto"
)

type CategoryRepository interface {
	ListAll(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, params repository.CreateCategoryParams) (int, error)
}

type CategoryService struct {
	categoryRepository CategoryRepository
}

func NewCategoryService(categoryRepository CategoryRepository) CategoryService {
	return CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s CategoryService) Create(ctx context.Context, createCategoryDTO dto.CreateCategoryDTO) (int, error) {
	categoryID, err := s.categoryRepository.Create(ctx, repository.CreateCategoryParams{Name: createCategoryDTO.Name, Slug: createCategoryDTO.Slug})
	if err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (s CategoryService) ListAll(ctx context.Context) ([]dto.GetCategoryDTO, error) {
	categories, err := s.categoryRepository.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return []dto.GetCategoryDTO{}, nil
	}

	dtos := make([]dto.GetCategoryDTO, len(categories))
	for i := range dtos {
		dtos[i] = dto.ToGetCategoryDTO(categories[i])
	}

	return dtos, nil
}
