package model

import (
	"fmt"
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"github.com/vlaship/book-catalog-go/internal/apperr"
	"github.com/vlaship/book-catalog-go/pkg/utils/mask"
)

// User is a model for user
type User struct {
	ID       types.UserID   `db:"user_id"`
	Username types.Username `db:"username"`
	Password types.Password `db:"password"`
	Data     UserData       `jsonb:"user_data"`
}

// String
func (u *User) String() string {
	return fmt.Sprintf(
		"User{ID: %d, Username: %s, Data: %s}",
		u.ID,
		mask.String(string(u.Username)),
		u.Data,
	)
}

// UserData is a model for user info
type UserData struct {
	FirstName string     `json:"firstname,omitempty"`
	LastName  string     `json:"lastname,omitempty"`
	Email     string     `json:"email,omitempty"`
	Plan      string     `json:"user_plan,omitempty"`
	Status    UserStatus `json:"status,omitempty"`
}

// String
func (ui *UserData) String() string {
	return fmt.Sprintf(
		"UserData{FirstName: %s, LastName: %s, Email: %s}",
		mask.String(ui.FirstName),
		mask.String(ui.LastName),
		mask.String(ui.Email),
	)
}

// UserStatus codes
type UserStatus string

// UserStatus codes
const (
	UserStatusActive       UserStatus = ""
	UserStatusNotActivated UserStatus = "user_status_not_activated"
)

// GetAppError by user status
func (u *User) GetAppError() error {
	switch u.Data.Status {
	case UserStatusActive:
		return nil
	case UserStatusNotActivated:
		return apperr.ErrUserNotActivated
	default:
		return apperr.ErrForbidden
	}
}
