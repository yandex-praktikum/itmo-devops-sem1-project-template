package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoadResult struct {
	TotalQuantity   int             `json:"total_items"` //nolint:tagliatelle //неправильно обрабатывает
	TotalCategories int             `json:"total_categories"`
	TotalPrice      decimal.Decimal `json:"total_price"`
}

type Product struct {
	ID         int             `db:"id"`
	Name       string          `db:"name"`
	Category   string          `db:"category"`
	Price      decimal.Decimal `db:"price"`
	CreateDate time.Time       `db:"create_date"`
}
