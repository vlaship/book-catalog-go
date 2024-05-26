package controller

import (
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/httphandling"
	"book-catalog/internal/logger"
	"book-catalog/internal/validation"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// UserReader is an interface for user reader
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-user-reader.go -package=mock . UserReader
type UserReader interface {
	GetUser(ctx context.Context) response.User
}

// UserWriter is an interface for user writer
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-user-writer.go -package=mock . UserWriter
type UserWriter interface {
	UpdateInfo(ctx context.Context, req *request.UserData) error
}

// UserController is a controller for user
type UserController struct {
	reader UserReader
	writer UserWriter
	valid  validation.Validator
	eh     httphandling.HTTPErrorHandler
	log    logger.Logger
}

// NewUserController creates a new UserController instance.
func NewUserController(
	reader UserReader,
	writer UserWriter,
	valid validation.Validator,
	eh httphandling.HTTPErrorHandler,
	log logger.Logger,
) *UserController {
	return &UserController{
		reader: reader,
		writer: writer,
		valid:  valid,
		eh:     eh,
		log:    log.New("UserController"),
	}
}

// RegisterRoutes registers routes for user controller
func (ctrl *UserController) RegisterRoutes(router chi.Router) {
	ctrl.log.Trc().Msg("RegisterRoutes")

	router.Route("/user", func(r chi.Router) {
		r.Get("/", ctrl.eh.HandlerError(ctrl.GetUser))
		r.Put("/info", ctrl.eh.HandlerError(ctrl.UpdateInfo))
	})
}

// GetUser returns user
// @Tags User
// @Security BearerAuth
// @Produce      json
// @Success 200 {object} response.User
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /user [get]
func (ctrl *UserController) GetUser(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("GetUser")

	resp := ctrl.reader.GetUser(r.Context())

	return encode(w, resp)
}

// UpdateInfo updates user info
// @Tags User
// @Security BearerAuth
// @Accept  json
// @Param user body request.UserData true "User"
// @Success 200 "OK" "Accepted"
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /user/info [put]
func (ctrl *UserController) UpdateInfo(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("UpdateInfo")

	req, err := ctrl.unmarshal(w, r)
	if err != nil {
		return err
	}

	if err = ctrl.writer.UpdateInfo(r.Context(), req); err != nil {
		return addTitle(err, "Problem updating user info")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func (ctrl *UserController) unmarshal(w http.ResponseWriter, r *http.Request) (*request.UserData, error) {
	ctrl.log.Trc().Ctx(r.Context()).Msg("unmarshal")

	var req = request.UserData{}
	res, err := decode(w, r, &req)
	if err != nil {
		ctrl.log.Err(err).Msg(unmarshalRequest)
		return nil, err
	}

	if err = validate(req, ctrl.valid); err != nil {
		return nil, err
	}

	return res, nil
}
