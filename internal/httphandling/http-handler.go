package httphandling

import (
	"book-catalog/internal/apperr"
	"net/http"
)

// HTTPErrorHandler interface
//
//go:generate mockgen -destination=../../test/mock/httphandling/mock-http_handler.go -package=mock . HTTPErrorHandler
type HTTPErrorHandler interface {
	HandlerError(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc

	AppErrorResponse(w http.ResponseWriter, r *http.Request, a apperr.AppError)
}
