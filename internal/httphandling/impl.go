package httphandling

import (
	"book-catalog/internal/app/dto/response"
	"book-catalog/internal/apperr"
	"book-catalog/internal/logger"
	"errors"
	"net/http"
	"time"
)

const (
	problemJSON       = "application/problem+json"
	headerContentType = "Content-Type"
)

// HTTPErrorHandlerImpl is a implementation of HTTPErrorHandler
type HTTPErrorHandlerImpl struct {
	log logger.Logger
}

// New returns a new HTTPErrorHandlerImpl
func New(log logger.Logger) HTTPErrorHandler {
	return &HTTPErrorHandlerImpl{log: log}
}

// HandlerError is a middleware that handles errors returned by HTTP handlers
func (h *HTTPErrorHandlerImpl) HandlerError(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}

		h.log.Wrn().Err(err).Values(
			"requestID", r.Header.Get("X-Request-ID"),
		).Msg("error handling request")

		p := getProblemDetails(err)
		write(w, r, p)
	}
}

func AppErrorResponse(w http.ResponseWriter, r *http.Request, a apperr.AppError) {
	p := getProblemDetails(a)
	write(w, r, p)
}

func write(w http.ResponseWriter, r *http.Request, p *response.ProblemDetail) {
	p.Timestamp = time.Now().Format(time.RFC3339Nano)
	p.Instance = r.RequestURI

	// Write the JSON response
	w.Header().Set(headerContentType, problemJSON)
	w.WriteHeader(p.Status)
	_, _ = w.Write(p.JSON())
}

func newFromError(err error) response.ProblemDetail {
	return response.ProblemDetail{
		Detail: err.Error(),
	}
}

func newFromAppError(err apperr.AppError) response.ProblemDetail {
	return response.ProblemDetail{
		Title:  err.Title,
		Detail: err.Detail,
		Code:   err.Code,
	}
}

func getProblemDetails(err error) *response.ProblemDetail {
	var appError apperr.AppError
	if errors.As(err, &appError) {
		p := newFromAppError(appError)
		p.Status = getStatus(appError)
		return &p
	}

	p := newFromError(err)
	p.Status = http.StatusInternalServerError
	p.Title = http.StatusText(http.StatusInternalServerError)
	return &p
}

func getStatus(err error) int {
	switch {
	case errors.Is(err, apperr.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, apperr.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperr.ErrForbidden),
		errors.Is(err, apperr.ErrInvalidToken),
		errors.Is(err, apperr.ErrUserNotActivated),
		errors.Is(err, apperr.ErrInvalidOTP):
		return http.StatusForbidden
	case errors.Is(err, apperr.ErrBadRequest),
		errors.Is(err, apperr.ErrAlreadyExists):
		return http.StatusBadRequest
	case errors.Is(err, apperr.ErrUnsupportedMediaType):
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}
