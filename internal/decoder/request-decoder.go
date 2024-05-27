package decoder

import (
	"book-catalog/internal/apperr"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

const (
	megabyte          = 1048576
	headerContentType = "Content-Type"
	applicationJSON   = "application/json"
)

// Decode is a helper function to decode JSON requests
func Decode(w http.ResponseWriter, r *http.Request, dst any) error { //nolint:cyclop // it's ok
	if r.Header.Get(headerContentType) != "" {
		value, _ := header.ParseValueAndParams(r.Header, headerContentType)
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
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return apperr.ErrDecodingRequest.WithFunc(
			apperr.WithDetail(fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)))

	case errors.Is(err, io.ErrUnexpectedEOF):
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body contains badly-formed JSON"))

	case errors.As(err, &unmarshalTypeError):
		return apperr.ErrDecodingRequest.WithFunc(
			apperr.WithDetail(fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset,
			),
			))

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail(fmt.Sprintf("Request body contains unknown field %s", fieldName)))

	case errors.Is(err, io.EOF):
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body must not be empty"))

	case err.Error() == "http: request body too large":
		return apperr.ErrDecodingRequest.WithFunc(apperr.WithDetail("Request body must not be larger than 1MB"))

	default:
		return err
	}
}
