package service

import (
	"context"
	"fmt"
	"log"
	"myapp/internal/domain"
	"myapp/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	UserRepository      repository.UserRepository
	UserRoleRepository  repository.UserRoleRepository
	CartRepository      repository.CartRepository
	CartItemRepository  repository.CartItemRepository
	ProductRepository   repository.ProductRepository
	OrderRepository     repository.OrderRepository
	OrderItemRepository repository.OrderItemRepository
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
		return nil, fmt.Errorf("transaction start: %w", domain.ErrInternalServerError)
	}

	repos := Repositories{
		UserRepository:      repository.NewUserRepository(tx),
		UserRoleRepository:  repository.NewUserRoleRepository(tx),
		CartRepository:      repository.NewCartRepository(tx),
		ProductRepository:   repository.NewProductRepository(tx),
		CartItemRepository:  repository.NewCartItemRepository(tx),
		OrderRepository:     repository.NewOrderRepository(tx),
		OrderItemRepository: repository.NewOrderItemRepository(tx),
	}
	value, err := fn(ctx, repos)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", domain.ErrInternalServerError)
	}

	return value, nil
}
