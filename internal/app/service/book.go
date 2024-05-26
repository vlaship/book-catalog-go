package service

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/logger"
	"context"
)

// BookReader is an interface for book reader
//
//go:generate mockgen -destination=../../../test/mock/service/mock-book-reader.go -package=mock . BookReader
type BookReader interface {
	GetBook(ctx context.Context, bookID types.ID) (*model.Book, error)
	GetBooks(ctx context.Context) ([]model.Book, error)
}

// BookWriter is an interface for book writer
//
//go:generate mockgen -destination=../../../test/mock/service/mock-book-writer.go -package=mock . BookWriter
type BookWriter interface {
	CreateBook(ctx context.Context, book *model.Book) (*model.Book, error)
	UpdateBook(ctx context.Context, bookID types.ID, book *model.Book) error
	DeleteBook(ctx context.Context, bookID types.ID) error
}

// BookService is a service for book
type BookService struct {
	reader BookReader
	writer BookWriter
	log    logger.Logger
}

// NewBookService creates new book service
func NewBookService(
	reader BookReader,
	writer BookWriter,
	log logger.Logger,
) *BookService {
	return &BookService{
		reader: reader,
		writer: writer,
		log:    log.New("BookService"),
	}
}

// GetBook returns book by id
func (s *BookService) GetBook(ctx context.Context, bookID types.ID) (*model.Book, error) {
	s.log.Dbg().Ctx(ctx).Values("bookID", bookID).Msg("GetBook")

	return s.reader.GetBook(ctx, bookID)
}

// GetBooks returns all books
func (s *BookService) GetBooks(ctx context.Context) ([]model.Book, error) {
	s.log.Trc().Ctx(ctx).Msg("GetBooks")

	return s.reader.GetBooks(ctx)
}

// CreateBook creates new book
func (s *BookService) CreateBook(ctx context.Context, book *model.Book) (*model.Book, error) {
	s.log.Dbg().Ctx(ctx).Values("book", book).Msg("CreateBook")

	return s.writer.CreateBook(ctx, book)
}

// UpdateBook updates book
func (s *BookService) UpdateBook(ctx context.Context, bookID types.ID, book *model.Book) error {
	s.log.Dbg().Ctx(ctx).Values("bookID", bookID, "book", book).Msg("UpdateBook")

	return s.writer.UpdateBook(ctx, bookID, book)
}

// DeleteBook deletes book
func (s *BookService) DeleteBook(ctx context.Context, bookID types.ID) error {
	s.log.Dbg().Ctx(ctx).Values("bookID", bookID).Msg("DeleteBook")

	return s.writer.DeleteBook(ctx, bookID)
}
