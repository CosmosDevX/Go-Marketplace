package dto

import "myapp/internal/domain"

type CreateCategoryDTO struct {
	Name string `json:"category_name" validate:"required,min=5,max=50"`
	Slug string `json:"category_slug" validate:"required,min=5,max=60"`
}

type GetCategoryDTO struct {
	ID   int    `json:"category_id"`
	Name string `json:"category_name"`
	Slug string `json:"category_slug"`
}

func ToGetCategoryDTO(category domain.Category) GetCategoryDTO {
	return GetCategoryDTO{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}
