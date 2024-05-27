package response

import (
	"book-catalog/internal/app/types"
)

// User is a response model for user
type User struct {
	Username types.Username `json:"username" example:"email@email.com"`
	Info     UserInfo       `json:"info"`
}

// UserInfo is a response model for user info
type UserInfo struct {
	FirstName string `json:"firstname" example:"John"`
	LastName  string `json:"lastname" example:"Doe"`
	Email     string `json:"email" example:"email@email.com"`
}
