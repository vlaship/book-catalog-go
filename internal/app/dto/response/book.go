package response

import "book-catalog/internal/app/types"

// CreateBook response
type CreateBook struct {
	ID types.ID `json:"id"`
}

// Book response
type Book struct {
	ID          types.ID      `json:"id" example:"1"`
	Title       string        `json:"title" example:"Book Title"`
	Description string        `json:"description" example:"Book Description"`
	ISBN        string        `json:"isbn" example:"1234567890"`
	AuthorID    types.ID      `json:"author_id" example:"1"`
	Price       types.Decimal `json:"price" example:"15.99"`
}

// ListBook response
type ListBook struct {
	ID    types.ID `json:"id" example:"1"`
	Title string   `json:"title" example:"Book Title"`
}
