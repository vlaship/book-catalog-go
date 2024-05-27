package controller

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/app/types"
	"book-catalog/internal/httphandling"
	"book-catalog/internal/logger"
	"book-catalog/internal/validation"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const bookPath = "/v1/book"

// BookReader is an interface for book reader
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-book-reader.go -package=mock . BookReader
type BookReader interface {
	GetBook(ctx context.Context, bookID types.ID) (*response.Book, error)
	GetBooks(ctx context.Context) ([]response.ListBook, error)
}

// BookWriter is an interface for book writer
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-book-writer.go -package=mock . BookWriter
type BookWriter interface {
	CreateBook(ctx context.Context, req *request.CreateBook) (*response.CreateBook, error)
	UpdateBook(ctx context.Context, bookID types.ID, req *request.UpdateBook) error
	DeleteBook(ctx context.Context, bookID types.ID) error
}

// BookController is a controller for book
type BookController struct {
	reader  BookReader
	writer  BookWriter
	valid   validation.Validator
	handler httphandling.HTTPErrorHandler
	log     logger.Logger
}

// NewBookController creates new book controller
func NewBookController(
	reader BookReader,
	writer BookWriter,
	valid validation.Validator,
	handler httphandling.HTTPErrorHandler,
	log logger.Logger,
) *BookController {
	return &BookController{
		reader:  reader,
		writer:  writer,
		valid:   valid,
		handler: handler,
		log:     log.New("BookController"),
	}
}

// RegisterRoutes registers routes
func (ctrl *BookController) RegisterRoutes(router chi.Router) {
	ctrl.log.Trc().Msg("RegisterRoutes")

	router.Route(bookPath, func(r chi.Router) {
		r.Get("/", ctrl.handler.HandlerError(ctrl.GetBooks))
		r.Post("/", ctrl.handler.HandlerError(ctrl.CreateBook))

		r.Route("/{bookID}", func(r chi.Router) {
			r.Get("/", ctrl.handler.HandlerError(ctrl.GetBook))
			r.Put("/", ctrl.handler.HandlerError(ctrl.UpdateBook))
			r.Delete("/", ctrl.handler.HandlerError(ctrl.DeleteBook))
		})
	})
}

// GetBooks gets books
// @Summary Get books
// @Tags Books
// @Security BearerAuth
// @Produce      json
// @Success 200 {array} response.ListBook
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/book [get]
func (ctrl *BookController) GetBooks(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("GetBooks")

	res, err := ctrl.reader.GetBooks(r.Context())
	if err != nil {
		return addTitle(err, "Problem getting books")
	}

	return encode(w, res)
}

// GetBook gets book by id
// @Summary Get book by id
// @Tags Books
// @Security BearerAuth
// @Produce      json
// @Param bookID path int true "Book ID"
// @Success 200 {object} response.Book
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/book/{bookID} [get]
func (ctrl *BookController) GetBook(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("GetBook")

	bookID, err := getBookID(r)
	if err != nil {
		return err
	}

	res, err := ctrl.reader.GetBook(r.Context(), bookID)
	if err != nil {
		return addTitle(err, "Problem getting book")
	}

	return encode(w, res)

}

// CreateBook creates a new book
// @Summary Create a new book
// @Tags Books
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param book body request.CreateBook true "Book"
// @Success 200 {object} response.CreateBook
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/book [post]
func (ctrl *BookController) CreateBook(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("CreateBook")

	req, err := decode(w, r, &request.CreateBook{}, ctrl.valid)
	if err != nil {
		return err
	}

	resp, err := ctrl.writer.CreateBook(r.Context(), req)
	if err != nil {
		return addTitle(err, "Problem creating book")
	}

	return encode(w, resp)
}

// UpdateBook updates a book
// @Summary Update a book
// @Tags Books
// @Security BearerAuth
// @Accept  json
// @Param bookID path int true "Book ID"
// @Param book body request.UpdateBook true "Book"
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/book/{bookID} [put]
func (ctrl *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("UpdateBook")

	bookID, err := getBookID(r)
	if err != nil {
		return err
	}

	req, err := decode(w, r, &request.UpdateBook{}, ctrl.valid)
	if err != nil {
		return err
	}

	err = ctrl.writer.UpdateBook(r.Context(), bookID, req)
	if err != nil {
		return addTitle(err, "Problem updating book")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// DeleteBook deletes a book
// @Summary Delete a book
// @Tags Books
// @Security BearerAuth
// @Param bookID path int true "Book ID"
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/book/{bookID} [delete]
func (ctrl *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("DeleteBook")

	bookID, err := getBookID(r)
	if err != nil {
		return err
	}

	err = ctrl.writer.DeleteBook(r.Context(), bookID)
	if err != nil {
		return addTitle(err, "Problem deleting book")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
