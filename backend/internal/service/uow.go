package service

import (
	"context"
	"log"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	CategoryRepository repository.CategoryRepository
	ProductRepository  repository.ProductRepository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) (any, *domain.DomainError)) (any, *domain.DomainError)
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return unitOfWork{
		db: db,
	}
}

func (u unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) (any, *domain.DomainError)) (any, *domain.DomainError) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, &domain.DomainError{Code: constants.TransactionError, Message: "transaction start failed"}
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
		return nil, &domain.DomainError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return value, nil
}
