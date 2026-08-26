package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"sort"

	"github.com/shopspring/decimal"
)

type orderRow struct {
	ID          int             `db:"order_id"`
	UserID      int             `db:"user_id"`
	OrderStatus string          `db:"order_status_name"`
	OrderTotal  decimal.Decimal `db:"order_total"`
}

type orderWithItemRow struct {
	OrderID     int             `db:"order_id"`
	UserID      int             `db:"user_id"`
	OrderStatus string          `db:"order_status_name"`
	OrderTotal  decimal.Decimal `db:"order_total"`

	OrderItemID       int             `db:"order_item_id"`
	OrderItemQuantity int             `db:"order_item_quantity"`
	OrderItemTotal    decimal.Decimal `db:"order_item_total"`

	ProductID       int             `db:"product_id"`
	ProductName     string          `db:"product_name"`
	ProductDesc     string          `db:"product_description"`
	ProductPrice    decimal.Decimal `db:"product_price"`
	ProductImage    string          `db:"product_image"`
	ProductSellerID int             `db:"product_seller_id"`

	CategoryID   int    `db:"category_id"`
	CategoryName string `db:"category_name"`
	CategorySlug string `db:"category_slug"`
}

type OrderRepository struct {
	db DBTX
}

func NewOrderRepository(db DBTX) OrderRepository {
	return OrderRepository{
		db: db,
	}
}

func (r OrderRepository) Create(ctx context.Context, userID int, orderStatus string, orderTotal decimal.Decimal) (int, error) {
	query := `
		INSERT INTO orders(user_id, order_status_id, order_total)
		SELECT $1, (SELECT order_status_id FROM order_statuses WHERE order_status_name = $2), $3
		RETURNING order_id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, userID, orderStatus, orderTotal).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, fmt.Errorf("create order with userID %d: %w", userID, domain.ErrTimeout)
		}

		return 0, fmt.Errorf("create order with userID %d: %w", userID, err)
	}

	return id, nil
}

func (r OrderRepository) ListByUserID(ctx context.Context, userID int) ([]domain.Order, error) {
	query := `SELECT 
		o.order_id, o.user_id, o.order_total, os.order_status_name, oi.order_item_id, oi.order_item_total, oi.order_item_quantity,
		oi.product_name, oi.product_price, oi.product_seller_id,
		COALESCE(p.product_id, 0) AS product_id, COALESCE(p.product_description, '') AS product_description, COALESCE(p.product_image, '') AS product_image,
		COALESCE(c.category_id, 0) AS category_id, COALESCE(c.category_name, '') AS category_name, COALESCE(c.category_slug, '') AS category_slug
		FROM orders AS o
		INNER JOIN order_statuses AS os ON o.order_status_id = os.order_status_id
		INNER JOIN order_items AS oi ON oi.order_id = o.order_id
		LEFT JOIN products AS p ON p.product_id = oi.product_id
		LEFT JOIN categories AS c ON p.product_category_id = c.category_id
		WHERE o.user_id = $1
		ORDER BY o.order_id ASC
	`

	var rows []orderWithItemRow
	if err := r.db.SelectContext(ctx, &rows, query, userID); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("list orders by userID %d: %w", userID, domain.ErrTimeout)
		}

		return nil, fmt.Errorf("list orders by userID %d: %w", userID, err)
	}

	orderItemsMap := map[int][]domain.OrderItem{}
	for i := range rows {
		orderItemsMap[rows[i].OrderID] = append(orderItemsMap[rows[i].OrderID], domain.OrderItem{
			ID:      rows[i].OrderItemID,
			OrderID: rows[i].OrderID,
			Product: domain.Product{
				ID:          rows[i].ProductID,
				Name:        rows[i].ProductName,
				Description: rows[i].ProductDesc,
				Price:       rows[i].ProductPrice,
				Image:       rows[i].ProductImage,
				Category: domain.Category{
					ID:   rows[i].CategoryID,
					Name: rows[i].CategoryName,
					Slug: rows[i].CategorySlug,
				},
				SellerID: rows[i].ProductSellerID,
			},
			Quantity: rows[i].OrderItemQuantity,
			Total:    rows[i].OrderItemTotal,
		})
	}

	ordersMap := map[int]domain.Order{}
	for i := range rows {
		if _, exists := ordersMap[rows[i].OrderID]; !exists {
			ordersMap[rows[i].OrderID] = domain.Order{
				ID:         rows[i].OrderID,
				UserID:     rows[i].UserID,
				Status:     rows[i].OrderStatus,
				OrderItems: orderItemsMap[rows[i].OrderID],
				Total:      rows[i].OrderTotal,
			}
		}
	}

	var domainModels []domain.Order
	for _, order := range ordersMap {
		domainModels = append(domainModels, order)
	}

	sort.Slice(domainModels, func(i int, j int) bool {
		return domainModels[i].ID < domainModels[j].ID
	})

	return domainModels, nil
}
