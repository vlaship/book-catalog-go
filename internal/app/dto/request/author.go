package request

import "github.com/vlaship/book-catalog-go/internal/app/types"

// CreateAuthor request
type CreateAuthor struct {
	Name string        `json:"name" validate:"required,min=1,max=255"`
	Dob  types.DateDay `json:"dob" validate:"required"`
}

// UpdateAuthor request
type UpdateAuthor struct {
	Name string        `json:"name" validate:"required,min=1,max=255"`
	Dob  types.DateDay `json:"dob" validate:"required" swaggertype:"primitive,string" example:"2021-01-01"`
}
