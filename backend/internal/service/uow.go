package service

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/logger"
	"myapp/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type userRepository interface {
	Create(ctx context.Context, u domain.User) (int, error)
}

type userRoleRepository interface {
	Create(ctx context.Context, username string, roleName string) error
}

type cartRepository interface {
	Create(ctx context.Context, userID int) (int, error)
	GetIDByUserID(ctx context.Context, userID int) (int, error)
}

type cartItemRepository interface {
	DeleteAllByCartID(ctx context.Context, cartID int) error
	ListByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error)
}

type orderRepository interface {
	Create(ctx context.Context, userID int, orderStatus string, orderTotal decimal.Decimal) (int, error)
}

type orderItemRepository interface {
	CreateMany(ctx context.Context, orderItems []domain.OrderItem) error
}

type Repositories struct {
	UserRepository      userRepository
	UserRoleRepository  userRoleRepository
	CartRepository      cartRepository
	CartItemRepository  cartItemRepository
	OrderRepository     orderRepository
	OrderItemRepository orderItemRepository
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
		CartItemRepository:  repository.NewCartItemRepository(tx),
		OrderRepository:     repository.NewOrderRepository(tx),
		OrderItemRepository: repository.NewOrderItemRepository(tx),
	}
	value, err := fn(ctx, repos)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.FromContext(ctx).Error("transaction rollback failed", "error", err)
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", domain.ErrInternalServerError)
	}

	return value, nil
}
