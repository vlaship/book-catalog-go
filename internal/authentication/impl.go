package authentication

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"github.com/vlaship/book-catalog-go/internal/apperr"
	"github.com/vlaship/book-catalog-go/internal/config"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	bearer     = "Bearer "
	authHeader = "Authorization"
)

// AuthenticatorImpl is an implementation of the Authenticator interface.
type AuthenticatorImpl struct {
	secret   []byte
	duration time.Duration
	method   jwt.SigningMethod
}

// New creates a new auth instance.
func New(cfg *config.Config) Authenticator {
	return &AuthenticatorImpl{
		secret:   cfg.JWT.Secret,
		duration: cfg.JWT.Duration,
		method:   jwt.SigningMethodHS256,
	}
}

// GenerateAccessToken creates a new access token.
func (j *AuthenticatorImpl) GenerateAccessToken(userID types.UserID) (accessToken types.Token, expiresIn int64, err error) {
	token := jwt.NewWithClaims(j.method, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
		// * add more claims
	})
	t, err := token.SignedString(j.secret)
	return types.Token(t), j.expiresIn(), err
}

// GetUserID returns the user id from the request.
func (j *AuthenticatorImpl) GetUserID(r *http.Request) (types.UserID, error) {
	userID, err := j.getUserID(r)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (j *AuthenticatorImpl) expiresIn() int64 {
	return j.duration.Milliseconds() / 1000
}

func (j *AuthenticatorImpl) getUserID(r *http.Request) (types.UserID, error) {
	t, err := j.extractToken(r)
	if err != nil {
		return 0, err
	}

	claims, err := j.extractClaims(t)
	if err != nil {
		return 0, err
	}

	return j.extractUserID(claims)
}

func (j *AuthenticatorImpl) extractToken(r *http.Request) (string, error) {
	h := r.Header.Get(authHeader)
	if !strings.HasPrefix(h, bearer) {
		return "", apperr.ErrNoBearerToken
	}

	return strings.TrimPrefix(h, bearer), nil
}

func (j *AuthenticatorImpl) extractClaims(token string) (jwt.Claims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if m := token.Method.Alg(); m != j.method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", m)
		}

		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if t.Valid {
		return t.Claims, nil
	}

	return nil, err
}

func (j *AuthenticatorImpl) extractUserID(claims jwt.Claims) (types.UserID, error) {
	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return 0, jwt.ErrTokenInvalidClaims
	}

	now := time.Now()
	if exp.Before(now) || exp.Equal(now) {
		return 0, jwt.ErrTokenExpired
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return 0, jwt.ErrTokenInvalidSubject
	}

	return types.UserID(id), nil
}
