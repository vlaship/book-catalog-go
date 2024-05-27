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

const authorPath = "/v1/author"

// AuthorReader is an interface for author reader
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-author-reader.go -package=mock . AuthorReader
type AuthorReader interface {
	GetAuthors(ctx context.Context) ([]response.ListAuthor, error)
	GetAuthor(ctx context.Context, authorID types.ID) (*response.Author, error)
}

// AuthorWriter is an interface for author writer
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-author-writer.go -package=mock . AuthorWriter
type AuthorWriter interface {
	CreateAuthor(ctx context.Context, req *request.CreateAuthor) (*response.CreateAuthor, error)
	UpdateAuthor(ctx context.Context, authorID types.ID, author *request.UpdateAuthor) error
	DeleteAuthor(ctx context.Context, authorID types.ID) error
}

// AuthorController is a controller for author
type AuthorController struct {
	reader  AuthorReader
	writer  AuthorWriter
	valid   validation.Validator
	handler httphandling.HTTPErrorHandler
	log     logger.Logger
}

// NewAuthorController creates new author controller
func NewAuthorController(
	reader AuthorReader,
	writer AuthorWriter,
	valid validation.Validator,
	handler httphandling.HTTPErrorHandler,
	log logger.Logger,
) *AuthorController {
	return &AuthorController{
		reader:  reader,
		writer:  writer,
		valid:   valid,
		handler: handler,
		log:     log.New("AuthorController"),
	}
}

// RegisterRoutes registers author routes
func (ctrl *AuthorController) RegisterRoutes(router chi.Router) {
	ctrl.log.Trc().Msg("RegisterRoutes")

	router.Route(authorPath, func(r chi.Router) {
		r.Get("/", ctrl.handler.HandlerError(ctrl.GetAuthors))
		r.Post("/", ctrl.handler.HandlerError(ctrl.CreateAuthor))

		r.Route("/{authorID}", func(r chi.Router) {
			r.Get("/", ctrl.handler.HandlerError(ctrl.GetAuthor))
			r.Put("/", ctrl.handler.HandlerError(ctrl.UpdateAuthor))
			r.Delete("/", ctrl.handler.HandlerError(ctrl.DeleteAuthor))
		})
	})
}

// GetAuthors gets authors
// @Summary Get authors
// @Tags Author
// @Security BearerAuth
// @Produce      json
// @Success 200 {array} response.ListAuthor
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/author [get]
func (ctrl *AuthorController) GetAuthors(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("GetAuthors")

	res, err := ctrl.reader.GetAuthors(r.Context())
	if err != nil {
		return addTitle(err, "Problem getting authors")
	}

	return encode(w, res)
}

// GetAuthor gets author by id
// @Summary Get author by id
// @Tags Author
// @Security BearerAuth
// @Produce      json
// @Param authorID path int true "Author ID"
// @Success 200 {object} response.Author
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/author/{authorID} [get]
func (ctrl *AuthorController) GetAuthor(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("GetAuthor")

	authorID, err := getAuthorID(r)
	if err != nil {
		return err
	}

	res, err := ctrl.reader.GetAuthor(r.Context(), authorID)
	if err != nil {
		return addTitle(err, "Problem getting author")
	}

	return encode(w, res)
}

// CreateAuthor creates author
// @Summary Create author
// @Tags Author
// @Security BearerAuth
// @Accept      json
// @Produce      json
// @Param author body request.CreateAuthor true "Author"
// @Success 200 {object} response.Author
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/author [post]
func (ctrl *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("CreateAuthorReq")

	req, err := ctrl.unmarshalCreate(w, r)
	if err != nil {
		return err
	}

	res, err := ctrl.writer.CreateAuthor(r.Context(), req)
	if err != nil {
		return addTitle(err, "Problem creating author")
	}

	w.WriteHeader(http.StatusOK)

	return encode(w, res)
}

// UpdateAuthor updates author
// @Summary Update author
// @Tags Author
// @Security BearerAuth
// @Accept      json
// @Param authorID path int true "Author ID"
// @Param author body request.UpdateAuthor true "Author"
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/author/{authorID} [put]
func (ctrl *AuthorController) UpdateAuthor(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("UpdateAuthorReq")

	authorID, err := getAuthorID(r)
	if err != nil {
		return err
	}
	req, err := ctrl.unmarshalUpdate(w, r)
	if err != nil {
		return err
	}

	err = ctrl.writer.UpdateAuthor(r.Context(), authorID, req)
	if err != nil {
		return addTitle(err, "Problem updating author")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// DeleteAuthor deletes author
// @Summary Delete author
// @Tags Author
// @Security BearerAuth
// @Param authorID path int true "Author ID"
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 404 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/author/{authorID} [delete]
func (ctrl *AuthorController) DeleteAuthor(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("DeleteAuthor")

	authorID, err := getAuthorID(r)
	if err != nil {
		return err
	}

	err = ctrl.writer.DeleteAuthor(r.Context(), authorID)
	if err != nil {
		return addTitle(err, "Problem deleting author")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func (ctrl *AuthorController) unmarshalCreate(
	w http.ResponseWriter,
	r *http.Request,
) (*request.CreateAuthor, error) {
	ctrl.log.Trc().Msg(unmarshalRequest)

	req, err := decode(w, r, &request.CreateAuthor{})
	if err != nil {
		ctrl.log.Err(err).Msg(unmarshalRequest)
		return nil, err
	}

	if err := validate(req, ctrl.valid); err != nil {
		return nil, err
	}

	return req, nil
}

func (ctrl *AuthorController) unmarshalUpdate(
	w http.ResponseWriter,
	r *http.Request,
) (*request.UpdateAuthor, error) {
	ctrl.log.Trc().Msg(unmarshalRequest)

	req, err := decode(w, r, &request.UpdateAuthor{})
	if err != nil {
		ctrl.log.Err(err).Msg(unmarshalRequest)
		return nil, err
	}

	if err = validate(req, ctrl.valid); err != nil {
		return nil, err
	}

	return req, nil
}
