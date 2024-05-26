package httphandling

import (
	"book-catalog/internal/apperr"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetStatus(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{apperr.ErrNotFound, http.StatusNotFound},
		{apperr.ErrBadRequest, http.StatusBadRequest},
		{apperr.ErrAlreadyExists, http.StatusBadRequest},
		{apperr.ErrUnsupportedMediaType, http.StatusUnsupportedMediaType},
		{apperr.ErrUnauthorized, http.StatusUnauthorized},
		{apperr.ErrForbidden, http.StatusForbidden},
		{apperr.ErrInvalidToken, http.StatusForbidden},
		{apperr.ErrUserNotActivated, http.StatusForbidden},
		{apperr.ErrInvalidOTP, http.StatusForbidden},
		{apperr.ErrExecuteTemplate, http.StatusInternalServerError},
		{errors.New(""), http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.err.Error(), func(t *testing.T) {
			result := getStatus(test.err)
			assert.Equal(t, test.expected, result, "Unexpected status code for error: %s", test.err)
		})
	}
}
