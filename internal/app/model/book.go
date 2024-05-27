package model

import (
	"book-catalog/internal/app/types"
	"github.com/shopspring/decimal"
)

// Book model
type Book struct {
	ID          types.ID        `db:"id"`
	Title       string          `db:"title"`
	Description string          `db:"description"`
	ISBN        string          `db:"isbn"`
	AuthorID    types.ID        `db:"author_id"`
	Price       decimal.Decimal `db:"price"`
}
