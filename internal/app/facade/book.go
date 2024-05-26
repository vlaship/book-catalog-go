package facade

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/mapper"
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/logger"
	"context"
)

// BookReader is an interface for book reader
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-book-reader.go -package=mock . BookReader
type BookReader interface {
	GetBook(ctx context.Context, bookID types.ID) (*model.Book, error)
	GetBooks(ctx context.Context) ([]model.Book, error)
}

// BookWriter is an interface for book writer
//
//go:generate mockgen -destination=../../../test/mock/facade/mock-book-writer.go -package=mock . BookWriter
type BookWriter interface {
	CreateBook(ctx context.Context, book *model.Book) (*model.Book, error)
	UpdateBook(ctx context.Context, bookID types.ID, book *model.Book) error
	DeleteBook(ctx context.Context, bookID types.ID) error
}

// BookFacade is a facade for book
type BookFacade struct {
	reader BookReader
	writer BookWriter
	m      mapper.Book
	log    logger.Logger
}

// NewBookFacade creates new book facade
func NewBookFacade(
	reader BookReader,
	writer BookWriter,
	log logger.Logger,
) *BookFacade {
	return &BookFacade{
		reader: reader,
		writer: writer,
		m:      mapper.Book{},
		log:    log.New("BookFacade"),
	}
}

// GetBook returns book by id
func (f *BookFacade) GetBook(ctx context.Context, bookID types.ID) (*response.Book, error) {
	f.log.Dbg().Ctx(ctx).Values("bookID", bookID).Msg("GetBook")

	book, err := f.reader.GetBook(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return f.m.BookResp(book), nil
}

// GetBooks returns all books
func (f *BookFacade) GetBooks(ctx context.Context) ([]response.ListBook, error) {
	f.log.Trc().Ctx(ctx).Msg("GetBooks")

	books, err := f.reader.GetBooks(ctx)
	if err != nil {
		return nil, err
	}

	return f.m.BooksResp(books), nil
}

// CreateBook creates new book
func (f *BookFacade) CreateBook(ctx context.Context, book *request.CreateBook) (*response.CreateBook, error) {
	f.log.Dbg().Ctx(ctx).Values("book", book).Msg("CreateBookReq")

	req := f.m.CreateBookReq(book)
	resp, err := f.writer.CreateBook(ctx, req)
	if err != nil {
		return nil, err
	}

	return f.m.CreateBookResp(resp), nil
}

// UpdateBook updates book by id
func (f *BookFacade) UpdateBook(ctx context.Context, bookID types.ID, book *request.UpdateBook) error {
	f.log.Dbg().Ctx(ctx).Values("book", book).Msg("UpdateBookReq")

	return f.writer.UpdateBook(ctx, bookID, f.m.UpdateBookReq(book))
}

// DeleteBook deletes book by id
func (f *BookFacade) DeleteBook(ctx context.Context, bookID types.ID) error {
	f.log.Dbg().Ctx(ctx).Values("bookID", bookID).Msg("DeleteBook")

	return f.writer.DeleteBook(ctx, bookID)
}
