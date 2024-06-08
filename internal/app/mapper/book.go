package mapper

import (
	"github.com/vlaship/book-catalog-go/internal/app/dto/request"
	"github.com/vlaship/book-catalog-go/internal/app/dto/response"
	"github.com/vlaship/book-catalog-go/internal/app/model"
	"github.com/vlaship/book-catalog-go/internal/app/types"
)

// Book is a mapper for book
type Book struct{}

// CreateBookReq creates a new book model
func (m *Book) CreateBookReq(req *request.CreateBook) *model.Book {
	return &model.Book{
		Title:       req.Title,
		Description: req.Description,
		ISBN:        req.ISBN,
		AuthorID:    req.AuthorID,
		Price:       req.Price.Value,
	}
}

// UpdateBookReq updates a book model
func (m *Book) UpdateBookReq(req *request.UpdateBook) *model.Book {
	return &model.Book{
		Title:       req.Title,
		Description: req.Description,
		ISBN:        req.ISBN,
		AuthorID:    req.AuthorID,
		Price:       req.Price.Value,
	}
}

// CreateBookResp creates a new book response
func (m *Book) CreateBookResp(out *model.Book) *response.CreateBook {
	return &response.CreateBook{
		ID: out.ID,
	}
}

// BookResp creates a new book response
func (m *Book) BookResp(out *model.Book) *response.Book {
	return &response.Book{
		ID:          out.ID,
		Title:       out.Title,
		Description: out.Description,
		ISBN:        out.ISBN,
		AuthorID:    out.AuthorID,
		Price:       types.Decimal{Decimal: out.Price},
	}
}

// BooksResp creates a new list of book response
func (m *Book) BooksResp(out []model.Book) []response.ListBook {
	books := make([]response.ListBook, 0, len(out))
	for i := range out {
		books = append(books, response.ListBook{
			ID:    out[i].ID,
			Title: out[i].Title,
		})
	}
	return books
}
