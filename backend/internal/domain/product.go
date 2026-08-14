package domain

import "github.com/shopspring/decimal"

type Product struct {
	ID          int             `db:"product_id"`
	Name        string          `db:"product_name"`
	Description string          `db:"product_description"`
	Price       decimal.Decimal `db:"product_price"`
	Quantity    int             `db:"product_quantity"`
	Image       string          `db:"product_image"`
	Category    Category        `db:"category"`
	CategoryID  int             `db:"product_category_id"`
}
