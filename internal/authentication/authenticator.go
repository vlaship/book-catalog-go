package authentication

import (
	"book-catalog/internal/app/types"
	"net/http"
)

// Authenticator interface
//
//go:generate mockgen -destination=../../test/mock/authentication/mock-authenticator.go -package=mock . Authenticator
type Authenticator interface {
	GenerateAccessToken(userID types.UserID) (accessToken types.Token, expiresIn int64, err error)
	GetUserID(r *http.Request) (types.UserID, error)
}
