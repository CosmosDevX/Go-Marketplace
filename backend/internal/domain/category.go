package domain

import "fmt"

type Category struct {
	ID   int
	Name string
	Slug string
}

func NewCategory(name, slug string) (Category, error) {
	if name == "" {
		return Category{}, fmt.Errorf("name: %w", ErrValidation)
	}
	if slug == "" {
		return Category{}, fmt.Errorf("slug: %w", ErrValidation)
	}

	return Category{
		Name: name,
		Slug: slug,
	}, nil
}
