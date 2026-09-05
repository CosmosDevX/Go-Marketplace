package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/lib/pq"
)

type cartRow struct {
	ID     int `db:"cart_id"`
	UserID int `db:"user_id"`
}

type CartRepository struct {
	db DBTX
}

func NewCartRepository(db DBTX) CartRepository {
	return CartRepository{
		db: db,
	}
}

func (r CartRepository) Create(ctx context.Context, userID int) (int, error) {
	query := `INSERT INTO carts(user_id) VALUES($1) RETURNING cart_id`
	var id int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create cart: %w", domain.ErrTimeout)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("create cart: %w", domain.ErrUniqueViolation)
		}

		return 0, fmt.Errorf("create cart: %w", err)
	}

	return id, nil
}

func (r CartRepository) GetIDByUserID(ctx context.Context, userID int) (int, error) {
	query := `SELECT cart_id FROM carts WHERE user_id = $1 FOR UPDATE`
	var id int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("get cartID by userID %d: %w", userID, domain.ErrTimeout)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("get cartID by userID %d: %w", userID, domain.ErrNotFound)
		}

		return 0, fmt.Errorf("get cartID by userID %d: %w", userID, err)
	}

	return id, nil
}
