package service

import (
	"context"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
)

type CategoryRepository interface {
	ListAll(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, c domain.Category) (int, error)
}

type CategoryService struct {
	categoryRepository CategoryRepository
}

func NewCategoryService(categoryRepository CategoryRepository) CategoryService {
	return CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s CategoryService) Create(ctx context.Context, dto dto.CreateCategoryDTO) (int, error) {
	category, err := domain.NewCategory(dto.Name, dto.Slug)
	if err != nil {
		return 0, err
	}

	categoryID, err := s.categoryRepository.Create(ctx, category)
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
