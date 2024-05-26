package repository

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/database"
	"book-catalog/internal/logger"
	"context"
)

// BookRepository is an interface for book repository
type BookRepository struct {
	pool database.ConnPool
	log  logger.Logger
}

const (
	entityNameBook = "book"
)

const (
	getBooks = `
	SELECT book_id, book_title, book_description, book_isbn, book_author_id, book_price
	FROM catalog.books
	WHERE deleted = FALSE;
`
	getBookByID = `
	SELECT book_id, book_title, book_description, book_isbn, book_author_id, book_price
	FROM catalog.books
	WHERE book_id = $1 AND deleted = FALSE;
`
	updateBookByID = `
	UPDATE catalog.books SET book_title = $2, book_description = $3, book_isbn = $4, book_author_id = $5, book_price = $6
	WHERE book_id = $1 AND deleted = FALSE;
`
	insertBook = `
	INSERT INTO catalog.books (book_title, book_description, book_isbn, book_author_id, book_price)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING book_id;
`
	deleteBookByID = `
	UPDATE catalog.books SET deleted = TRUE WHERE book_id = $1;
`
)

// NewBookRepository creates new book repository
func NewBookRepository(pool database.ConnPool, log logger.Logger) *BookRepository {
	return &BookRepository{
		pool: pool,
		log:  log.New("BookRepository"),
	}
}

func (r *BookRepository) l() logger.Logger {
	return r.log
}

func (r *BookRepository) p() database.ConnPool {
	return r.pool
}

// GetBooks get list of books
func (r *BookRepository) GetBooks(ctx context.Context) ([]model.Book, error) {
	r.log.Trc().Ctx(ctx).Msg("GetBooks")

	req := getEntity[model.Book]{
		query:      getBooks,
		entityName: entityNameBook,
		destinations: func(book *model.Book) []any {
			return []any{
				&book.ID,
				&book.Title,
				&book.Description,
				&book.ISBN,
				&book.AuthorID,
				&book.Price,
			}
		},
	}

	return getAll(ctx, r, req)
}

// GetBook get close book by user_id and book id
func (r *BookRepository) GetBook(ctx context.Context, bookID types.ID) (*model.Book, error) {
	r.log.Dbg().Ctx(ctx).Values("bookID", bookID).Msg("GetBook")

	req := getEntity[model.Book]{
		query:      getBookByID,
		entityName: entityNameBook,
		args:       []any{bookID},
		destinations: func(book *model.Book) []any {
			return []any{
				&book.ID,
				&book.Title,
				&book.Description,
				&book.ISBN,
				&book.AuthorID,
				&book.Price,
			}
		},
	}

	return getOne(ctx, r, req)
}

// CreateBook create book
func (r *BookRepository) CreateBook(ctx context.Context, book *model.Book) (*model.Book, error) {
	r.log.Dbg().Ctx(ctx).Values("book", book).Msg("CreateBook")

	req := createEntity[model.Book]{
		getEntity: getEntity[model.Book]{
			query:      insertBook,
			entityName: entityNameBook,
			args: []any{
				book.Title,
				book.Description,
				book.ISBN,
				book.AuthorID,
				book.Price,
			},
			destinations: func(book *model.Book) []any { return []any{&book.ID} }},
	}

	return create(ctx, r, req)
}

// UpdateBook update book by userID and ID
func (r *BookRepository) UpdateBook(
	ctx context.Context,
	bookID types.ID,
	book *model.Book,
) error {
	r.log.Dbg().Ctx(ctx).Values("bookID", bookID, "book", book).Msg("UpdateBook")

	req := execRequest{
		query:      updateBookByID,
		entityName: entityNameBook,
		args: []any{
			bookID,
			book.Title,
			book.Description,
			book.ISBN,
			book.AuthorID,
			book.Price,
		},
	}

	return exec(ctx, r, req)
}

// DeleteBook delete book by ID
func (r *BookRepository) DeleteBook(ctx context.Context, bookID types.ID) error {
	r.log.Dbg().Ctx(ctx).Values("bookID", bookID).Msg("DeleteBook")

	req := execRequest{
		query:      deleteBookByID,
		entityName: entityNameBook,
		args:       []any{bookID},
	}

	return exec(ctx, r, req)
}
