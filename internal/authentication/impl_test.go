package authentication

import (
	"book-catalog/internal/app/types"
	"book-catalog/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)

	// when
	token, _, err := auth.GenerateAccessToken(types.UserID(1))

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetUserIDSuccess(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)
	expected := types.UserID(1)

	// when
	token, _, _ := auth.GenerateAccessToken(expected)
	result, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + string(token)},
		},
	})

	// then
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetUserIDFailNoHeader(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "no bearer token found")
}

func TestGetUserIDFailNoBearer(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "no bearer token found")
}

func TestGetUserIDFailInvalidTokenWithSuffix(t *testing.T) { // given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)
	token, _, _ := auth.GenerateAccessToken(types.UserID(1))

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + string(token) + "invalid"},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token signature is invalid: signature is invalid")
}

func TestGetUserIDFailInvalidTokenAnotherSecret(t *testing.T) {
	// given
	auth1 := New(&config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	})
	auth2 := New(&config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("invalid"),
			Duration: 1 * time.Hour,
		},
	})
	expected := types.UserID(1)
	token, _, _ := auth1.GenerateAccessToken(expected)

	// when
	_, err := auth2.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + string(token)},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token signature is invalid: signature is invalid")
}

func TestGetUserIDFailInvalidTokenExpired(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: -1 * time.Hour,
		},
	}
	auth := New(cfg)
	expected := types.UserID(1)
	token, _, _ := auth.GenerateAccessToken(expected)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + string(token)},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token has invalid claims: token is expired")
}

func TestGetUserIDFailInvalidTokenNoSubject(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.Duration)),
	}).SignedString(cfg.JWT.Secret)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + token},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token has invalid subject")
}

func TestGetUserIDFailInvalidTokenNoExpirationTime(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: "1",
	}).SignedString(cfg.JWT.Secret)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + token},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token has invalid claims")
}

func TestGetUserIDFailInvalidTokenNoSigningMethod(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)
	token, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject:   "1",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + token},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token is unverifiable: error while executing keyfunc: unexpected signing method: none")
}

func TestGetUserIDFailInvalidTokenInvalidSigningMethod(t *testing.T) {
	// given
	cfg := &config.Config{
		JWT: struct {
			Secret   []byte
			Duration time.Duration
		}{
			Secret:   []byte("secret"),
			Duration: 1 * time.Hour,
		},
	}
	auth := New(cfg)
	token, _ := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.RegisteredClaims{
		Subject:   "1",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}).SignedString(cfg.JWT.Secret)

	// when
	_, err := auth.GetUserID(&http.Request{
		Header: map[string][]string{
			"Authorization": {bearer + token},
		},
	})

	// then
	assert.Error(t, err)
	assert.Error(t, err, "token is malformed: token contains an invalid number of segments")
}
