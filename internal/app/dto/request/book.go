package request

import "github.com/vlaship/book-catalog-go/internal/app/types"

// CreateBook request
type CreateBook struct {
	Title       string                `json:"title" validate:"required,min=1,max=255"`
	Description string                `json:"description" validate:"required,min=1,max=255"`
	ISBN        string                `json:"isbn" validate:"required,min=1,max=255"`
	AuthorID    types.ID              `json:"author_id" validate:"required"`
	Price       types.PositiveDecimal `json:"price" example:"15.99" swaggertype:"primitive,number" validate:"required"`
}

// UpdateBook request
type UpdateBook struct {
	Title       string                `json:"title" validate:"required,min=1,max=255"`
	Description string                `json:"description" validate:"required,min=1,max=255"`
	ISBN        string                `json:"isbn" validate:"required,min=1,max=255"`
	AuthorID    types.ID              `json:"author_id" validate:"required"`
	Price       types.PositiveDecimal `json:"price" example:"15.99" swaggertype:"primitive,number" validate:"required"`
}
