package model

import "book-catalog/internal/app/types"

// Signin model
type Signin struct {
	AccessToken  types.Token
	Type         string
	ExpiresIn    int64
	RefreshToken types.Token
}
