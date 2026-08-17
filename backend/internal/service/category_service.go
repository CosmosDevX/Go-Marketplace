package service

import (
	"context"
	"myapp/internal/domain"
)

type CreateCategoryInput struct {
	Name string
	Slug string
}

type CategoryRepository interface {
	List(ctx context.Context) ([]domain.Category, error)
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

func (s CategoryService) Create(ctx context.Context, input CreateCategoryInput) (int, error) {
	category, err := domain.NewCategory(input.Name, input.Slug)
	if err != nil {
		return 0, err
	}

	categoryID, err := s.categoryRepository.Create(ctx, category)
	if err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (s CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.categoryRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
