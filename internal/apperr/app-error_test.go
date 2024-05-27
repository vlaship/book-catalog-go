package apperr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithFunc(t *testing.T) {
	given := AppError{
		Title:  "title",
		Code:   "code",
		Detail: "detail",
	}
	result := given.WithFunc(WithTitle("new title"), WithDetail("new detail"))

	assert.Error(t, result)
	assert.Equal(t, "new title", result.Title)
	assert.Equal(t, "new detail", result.Detail)
	assert.Equal(t, given.Code, result.Code)
	assert.ErrorIs(t, result.Err, given)
}
