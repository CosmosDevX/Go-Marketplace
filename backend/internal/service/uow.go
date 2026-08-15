package service

import (
	"context"
	"fmt"
	"log"
	"myapp/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	CategoryRepository repository.CategoryRepository
	ProductRepository  repository.ProductRepository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) (any, error)) (any, error)
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return unitOfWork{
		db: db,
	}
}

func (u unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) (any, error)) (any, error) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("transaction start: %w", err)
	}

	repos := Repositories{
		CategoryRepository: repository.NewCategoryRepository(tx),
		ProductRepository:  repository.NewProductRepository(tx),
	}

	value, domainErr := fn(ctx, repos)
	if domainErr != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
		return nil, domainErr
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", err)
	}

	return value, nil
}
