package controller

import (
	"book-catalog/internal/app/dto"
	"book-catalog/internal/app/dto/request"
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/httphandling"
	"book-catalog/internal/logger"
	"book-catalog/internal/validation"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	authPath = "/v1/auth"
)

// Auth interface
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-auth.go -package=mock book-catalog/internal/app/controller Auth
type Auth interface {
	Signin(ctx context.Context, req *request.Signin) (*response.Signin, error)
	Signup(ctx context.Context, req *request.Signup) error
}

// PasswordResetHandler interface
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-password-reset-handler.go -package=mock book-catalog/internal/app/controller PasswordResetHandler
type PasswordResetHandler interface {
	Reset(ctx context.Context, req *request.ResetPassword) error
	Replace(ctx context.Context, req *request.ReplacePassword) error
}

// Activator interface
//
//go:generate mockgen -destination=../../../test/mock/controller/mock-activator.go -package=mock book-catalog/internal/app/controller Activator
type Activator interface {
	Activate(ctx context.Context, req *request.Activation) error
	Resend(ctx context.Context, req *request.ResendActivation) error
}

// AuthController is a controller for authentication.
type AuthController struct {
	auth  Auth
	act   Activator
	pass  PasswordResetHandler
	valid validation.Validator
	eh    httphandling.HTTPErrorHandler
	log   logger.Logger
}

// NewAuthController creates a new AuthController instance.
func NewAuthController(
	auth Auth,
	act Activator,
	pass PasswordResetHandler,
	valid validation.Validator,
	eh httphandling.HTTPErrorHandler,
	log logger.Logger,
) *AuthController {
	return &AuthController{
		auth:  auth,
		act:   act,
		pass:  pass,
		valid: valid,
		eh:    eh,
		log:   log.New("AuthController"),
	}
}

// RegisterRoutes registers authentication routes.
func (ctrl *AuthController) RegisterRoutes(router chi.Router) {
	ctrl.log.Trc().Msg("RegisterRoutes")

	router.Route(authPath, func(r chi.Router) {
		r.Post("/signin", ctrl.eh.HandlerError(ctrl.Signin))
		r.Post("/signup", ctrl.eh.HandlerError(ctrl.Signup))
		r.Post("/activation/activate", ctrl.eh.HandlerError(ctrl.Activate))
		r.Post("/activation/resend", ctrl.eh.HandlerError(ctrl.Resend))
		r.Post("/password/reset", ctrl.eh.HandlerError(ctrl.Reset))
		r.Post("/password/replace", ctrl.eh.HandlerError(ctrl.Replace))
	})
}

// Signin
// @Summary Signin
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param signin body request.Signin true "Signin"
// @Success 200 {object} response.Signin
// @Failure 400 {object} response.ProblemDetail
// @Failure 401 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/auth/signin [post]
func (ctrl *AuthController) Signin(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("Signin")

	req, err := decodeAuth(w, r, request.Signin{}, ctrl.valid)
	if err != nil {
		return err
	}

	res, err := ctrl.auth.Signin(r.Context(), req)
	if err != nil {
		return addTitle(err, "Problem signing in")
	}

	return encode(w, res)
}

// Signup
// @Summary Signup
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param signup body request.Signup true "Signup"
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/auth/signup [post]
func (ctrl *AuthController) Signup(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("Signup")

	req, err := decodeAuth(w, r, request.Signup{}, ctrl.valid)
	if err != nil {
		return err
	}

	if err = ctrl.auth.Signup(r.Context(), req); err != nil {
		return addTitle(err, "Problem signing up")
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

// Activate
// @Summary Activate Account
// @Tags Authentication
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/auth/activation/activate [post]
func (ctrl *AuthController) Activate(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("Activate")

	req, err := decodeAuth(w, r, request.Activation{}, ctrl.valid)
	if err != nil {
		return err
	}

	if err = ctrl.act.Activate(r.Context(), req); err != nil {
		return addTitle(err, "Problem activating user")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// Resend
// @Summary Resend Activation Mail
// @Tags Authentication
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/auth/activation/resend [post]
func (ctrl *AuthController) Resend(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("Resend")

	req, err := decodeAuth(w, r, request.ResendActivation{}, ctrl.valid)
	if err != nil {
		return err
	}

	if err = ctrl.act.Resend(r.Context(), req); err != nil {
		return addTitle(err, "Problem resending activation")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// Reset
// @Summary Reset Password
// @Tags Authentication
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/auth/password/reset [post]
func (ctrl *AuthController) Reset(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("Reset")

	req, err := decodeAuth(w, r, request.ResetPassword{}, ctrl.valid)
	if err != nil {
		return err
	}

	if err = ctrl.pass.Reset(r.Context(), req); err != nil {
		return addTitle(err, "Problem resetting password")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// Replace
// @Summary Replace Password
// @Tags Authentication
// @Success 200 "OK"
// @Failure 400 {object} response.ProblemDetail
// @Failure 403 {object} response.ProblemDetail
// @Failure 500 {object} response.ProblemDetail
// @Router /v1/auth/password/replace [post]
func (ctrl *AuthController) Replace(w http.ResponseWriter, r *http.Request) error {
	ctrl.log.Trc().Ctx(r.Context()).Msg("Replace")

	req, err := decodeAuth(w, r, request.ReplacePassword{}, ctrl.valid)
	if err != nil {
		return err
	}

	if err = ctrl.pass.Replace(r.Context(), req); err != nil {
		return addTitle(err, "Problem replacing password")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func decodeAuth[T dto.Auth](
	w http.ResponseWriter,
	r *http.Request,
	req T,
	valid validation.Validator,
) (*T, error) {
	res, err := decode(w, r, &req)
	if err != nil {
		return nil, err
	}

	if err = validate(req, valid); err != nil {
		return nil, err
	}

	return res, nil
}
