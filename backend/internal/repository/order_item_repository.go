package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type orderItemRow struct {
	ID              int             `db:"order_item_id"`
	OrderID         int             `db:"order_id"`
	ProductID       int             `db:"product_id"`
	ProductName     string          `db:"product_name"`
	ProductPrice    decimal.Decimal `db:"product_price"`
	ProductSellerID int             `db:"product_seller_id"`
	Quantity        int             `db:"order_item_quantity"`
	Total           decimal.Decimal `db:"order_item_total"`
}

type OrderItemRepository struct {
	db DBTX
}

func NewOrderItemRepository(db DBTX) OrderItemRepository {
	return OrderItemRepository{
		db: db,
	}
}

func (r OrderItemRepository) CreateMany(ctx context.Context, orderItems []domain.OrderItem) error {
	if len(orderItems) == 0 {
		return fmt.Errorf("order items is empty: %w", domain.ErrValidation)
	}

	var rows []orderItemRow
	for i := range orderItems {
		rows = append(rows, orderItemRow{
			OrderID:         orderItems[i].OrderID,
			ProductID:       orderItems[i].Product.ID,
			Total:           orderItems[i].Total,
			Quantity:        orderItems[i].Quantity,
			ProductName:     orderItems[i].Product.Name,
			ProductPrice:    orderItems[i].Product.Price,
			ProductSellerID: orderItems[i].Product.SellerID,
		})
	}

	query := `INSERT INTO order_items(order_id, product_id, order_item_total, order_item_quantity, product_name, product_price, product_seller_id) 
		VALUES(:order_id, :product_id, :order_item_total, :order_item_quantity, :product_name, :product_price, :product_seller_id)
	`
	sqlResult, err := r.db.NamedExecContext(ctx, query, rows)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("create order item in order %d: %w", rows[0].OrderID, domain.ErrTimeout)
		}

		return fmt.Errorf("create order item in order %d: %w", rows[0].OrderID, err)
	}

	if affectedRows, _ := sqlResult.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("create order item in order %d: %w", rows[0].OrderID, domain.ErrInternalServerError)
	}

	return nil
}
