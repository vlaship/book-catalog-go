package response

import "github.com/vlaship/book-catalog-go/internal/app/types"

// Signin response
type Signin struct {
	AccessToken  types.Token `json:"access_token"`
	Type         string      `json:"type" example:"Bearer"`
	ExpiresIn    int64       `json:"expires_in" example:"3600"`
	RefreshToken types.Token `json:"refresh_token"`
}

// Signup response
type Signup struct {
	Username types.Username `json:"username"`
}
