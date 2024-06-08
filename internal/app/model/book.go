package model

import (
	"github.com/shopspring/decimal"
	"github.com/vlaship/book-catalog-go/internal/app/types"
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
