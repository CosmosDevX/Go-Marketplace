package service

import (
	"context"
	"fmt"
	"log"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
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

	repos := Repositories{}
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
