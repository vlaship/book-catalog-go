package decoder

import (
	"book-catalog/internal/apperr"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDecodeValidJSON(t *testing.T) {
	reqBody := `{"key": "value"}`
	r := &http.Request{
		Header: http.Header{
			headerContentType: []string{applicationJSON},
		},
		Body: io.NopCloser(strings.NewReader(reqBody)),
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.NoError(t, err)
	assert.Equal(t, "value", dst.Key)
}

func TestDecodeInvalidJSON(t *testing.T) {
	reqBody := `{"key": "value"`
	r := &http.Request{
		Header: http.Header{
			headerContentType: []string{applicationJSON},
		},
		Body: io.NopCloser(strings.NewReader(reqBody)),
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.Error(t, err)
	assert.Truef(t, errors.Is(err, apperr.ErrDecodingRequest), "Expected ErrDecodingRequest, got %v", err)
	assert.Contains(t, err.Error(), "Request body contains badly-formed JSON")
}

func TestDecodeUnsupportedContentType(t *testing.T) {
	r := &http.Request{
		Header: http.Header{
			headerContentType: []string{"application/xml"},
		},
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.Error(t, err)
	assert.Truef(t, errors.Is(err, apperr.ErrUnsupportedMediaType), "Expected ErrUnsupportedMediaType, got %v", err)
	assert.Contains(t, err.Error(), "Content-Type header is not application/json")
}

func TestDecodeEmptyRequest(t *testing.T) {
	r := &http.Request{
		Body: io.NopCloser(strings.NewReader("")),
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.Error(t, err)
	assert.Truef(t, errors.Is(err, apperr.ErrDecodingRequest), "Expected ErrDecodingRequest, got %v", err)
	assert.Contains(t, err.Error(), "Request body must not be empty")
}

func TestDecodeRequestTooLarge(t *testing.T) {
	reqBody := `{"key": "` + strings.Repeat("a", megabyte) + `"}` // 1MB + 1 byte
	r := &http.Request{
		Header: http.Header{
			headerContentType: []string{applicationJSON},
		},
		Body: io.NopCloser(strings.NewReader(reqBody)),
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.Error(t, err)
	assert.Truef(t, errors.Is(err, apperr.ErrDecodingRequest), "Expected ErrDecodingRequest, got %v", err)
	assert.Contains(t, err.Error(), "Request body must not be larger than 1MB")
}

func TestDecodeUnknownField(t *testing.T) {
	reqBody := `{"key": "value", "unknownField": "test"}`
	r := &http.Request{
		Header: http.Header{
			headerContentType: []string{applicationJSON},
		},
		Body: io.NopCloser(strings.NewReader(reqBody)),
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.NoError(t, err, "Unknown fields are ignored")
}

func TestDecodeMultipleJSONObjects(t *testing.T) {
	reqBody := `{"key1": "value1"}{"key2": "value2"}`
	r := &http.Request{
		Header: http.Header{
			headerContentType: []string{applicationJSON},
		},
		Body: io.NopCloser(strings.NewReader(reqBody)),
	}
	var dst struct {
		Key string `json:"key"`
	}
	err := Decode(nil, r, &dst)
	assert.Error(t, err)
	assert.Truef(t, errors.Is(err, apperr.ErrDecodingRequest), "Expected ErrDecodingRequest, got %v", err)
	assert.Contains(t, err.Error(), "Request body must only contain a single JSON object")
}
