package response

import "book-catalog/internal/app/types"

// CreateAuthor response
type CreateAuthor struct {
	ID types.ID `json:"id"`
}

// Author response
type Author struct {
	ID   types.ID      `json:"id" example:"1"`
	Name string        `json:"name" example:"John Doe"`
	Dob  types.DateDay `json:"dob" swaggertype:"primitive,string" example:"2021-01-01"`
}

// ListAuthor response
type ListAuthor struct {
	ID   types.ID `json:"id" example:"1"`
	Name string   `json:"name" example:"John Doe"`
}
