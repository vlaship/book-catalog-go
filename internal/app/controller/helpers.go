package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vlaship/book-catalog-go/internal/app/dto"
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"github.com/vlaship/book-catalog-go/internal/apperr"
	"github.com/vlaship/book-catalog-go/internal/decoder"
	"github.com/vlaship/book-catalog-go/internal/validation"
	"net/http"
)

const (
	headerContentType = "Content-Type"
	applicationJSON   = "application/json"
	extractParam      = "Extract param"
)

// encode is a helper function to encode JSON responses
func encode(w http.ResponseWriter, resource any) error {
	// Convert the resource to JSON
	resp, err := json.Marshal(resource)
	if err != nil {
		return apperr.ErrBadRequest.WithFunc(apperr.WithTitle("Encode"))
	}

	// Set the Content-Type header to application/json
	w.Header().Set(headerContentType, applicationJSON)

	// Write the JSON response
	_, _ = w.Write(resp)

	return nil
}

// getBookID is a helper function to get bookID from request
func getBookID(r *http.Request) (types.ID, error) {
	param := chi.URLParam(r, "bookID")
	bookID, err := types.NewID(param)
	if err != nil {
		return 0, apperr.ErrBadRequest.WithFunc(
			apperr.WithDetail(fmt.Sprintf("invalid bookID %v", param)),
			apperr.WithTitle(extractParam),
		)
	}

	return bookID, nil
}

// getAuthorID is a helper function to get authorID from request
func getAuthorID(r *http.Request) (types.ID, error) {
	param := chi.URLParam(r, "authorID")
	authorID, err := types.NewID(param)
	if err != nil {
		return 0, apperr.ErrBadRequest.WithFunc(
			apperr.WithDetail(fmt.Sprintf("invalid authorID %v", param)),
			apperr.WithTitle(extractParam),
		)
	}

	return authorID, nil
}

// addTitle adds title to problem
func addTitle(err error, title string) error {
	var appError apperr.AppError
	if errors.As(err, &appError) {
		return appError.WithFunc(apperr.WithTitle(title))
	}

	return err
}

// decode is a helper function to decode JSON requests
func decode[T dto.Request](
	w http.ResponseWriter,
	r *http.Request,
	req *T,
	v validation.Validator,
) (*T, error) {
	if err := decoder.Decode(w, r, req); err != nil {
		return nil, err
	}
	if err := v.Struct(req); err != nil {
		return nil, apperr.ErrValidationRequest.WithFunc(apperr.WithDetail(err.Error()))
	}

	return req, nil
}
