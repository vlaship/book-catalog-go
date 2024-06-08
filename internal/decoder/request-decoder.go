package decoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vlaship/book-catalog-go/internal/apperr"
	"io"
	"net/http"
	"strings"
)

const (
	megabyte          = 1048576
	headerContentType = "Content-Type"
	applicationJSON   = "application/json"
)

// Decode is a helper function to decode JSON requests
func Decode(w http.ResponseWriter, r *http.Request, dst any) error { //nolint:cyclop // it's ok
	if r.Header.Get(headerContentType) != "" {
		// get content type fro header
		value := r.Header.Get(headerContentType)
		if value != applicationJSON {
			return apperr.ErrUnsupportedMediaType.WithFunc(apperr.WithDetail(headerContentType + " header is not " + applicationJSON))
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, megabyte)
	defer r.Body.Close()

	// init decoder
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&dst); err != nil {
		return handleDecodeError(err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body must only contain a single JSON object"))
	}

	return nil
}

func handleDecodeError(err error) error {
	var syntaxError *json.SyntaxError
	var typeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail(msg))

	case errors.Is(err, io.ErrUnexpectedEOF):
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body contains badly-formed JSON"))

	case errors.As(err, &typeError):
		msg := fmt.Sprintf(
			"Request body contains an invalid value for the %q field (at position %d)", typeError.Field, typeError.Offset)
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail(msg))

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail(msg))

	case errors.Is(err, io.EOF):
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body must not be empty"))

	case err.Error() == "http: request body too large":
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body must not be larger than 1MB"))

	default:
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail(err.Error()))
	}
}
