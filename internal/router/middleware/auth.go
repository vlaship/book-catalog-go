package middleware

import (
	"book-catalog/internal/app/model"
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/internal/httphandling"
	"book-catalog/internal/logger"
	"context"
	"errors"
	"net/http"
)

// UserReader is an interface for reading users from a database.
//
//go:generate mockgen -destination=../../../test/mock/middleware/mock-user-reader.go -package=mock . UserReader
type UserReader interface {
	GetUserByID(ctx context.Context, userID types.UserID) (*model.User, error)
}

// UserIDReader is an interface for reading user IDs from a request.
//
//go:generate mockgen -destination=../../../test/mock/middleware/mock-user-id-reader.go -package=mock . UserIDReader
type UserIDReader interface {
	GetUserID(r *http.Request) (types.UserID, error)
}

// AuthMiddleware is a middleware that validates JWT tokens.
type AuthMiddleware struct {
	userIDReader UserIDReader
	userReader   UserReader
	handler      httphandling.HTTPErrorHandler
	log          logger.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware instance.
func NewAuthMiddleware(
	userIDReader UserIDReader,
	userReader UserReader,
	handler httphandling.HTTPErrorHandler,
	log logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		userReader:   userReader,
		userIDReader: userIDReader,
		handler:      handler,
		log:          log.New("AuthMiddleware"),
	}
}

// Validation is a middleware that validates JWT tokens.
func (m *AuthMiddleware) Validation() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID, err := m.userIDReader.GetUserID(r)
			if err != nil {
				m.log.Wrn().Err(err).Ctx(ctx).Msg("failed to get user id")
				m.unauthorized(w, r)
				return
			}

			user, err := m.getUser(ctx, userID)
			if err != nil {
				m.unauthorized(w, r)
				return
			}

			if ok := m.validateUser(w, r, user); !ok {
				return
			}

			// Add the valid user to the request context
			ctx = context.WithValue(ctx, types.UserContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (m *AuthMiddleware) validateUser(w http.ResponseWriter, r *http.Request, user *model.User) (ok bool) {
	if user.Data.Status == "" {
		return true
	}

	err := user.GetAppError()
	switch {

	case errors.Is(err, apperr.ErrForbidden):
		m.handler.AppErrorResponse(w, r, apperr.ErrForbidden)

	case errors.Is(err, apperr.ErrUserNotActivated):
		m.handler.AppErrorResponse(w, r, apperr.ErrUserNotActivated)

	default:
		m.unauthorized(w, r)
	}

	return false
}

func (m *AuthMiddleware) unauthorized(w http.ResponseWriter, r *http.Request) {
	m.handler.AppErrorResponse(w, r, apperr.ErrUnauthorized)
}

func (m *AuthMiddleware) getUser(ctx context.Context, userID types.UserID) (*model.User, error) {
	user, err := m.userReader.GetUserByID(ctx, userID)
	if err != nil {
		m.log.Wrn().Err(err).Ctx(ctx).Msg("failed to get user")
		return nil, err
	}

	return user, nil
}
