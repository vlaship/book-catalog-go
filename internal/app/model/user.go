package model

import (
	"book-catalog/internal/app/types"
	"book-catalog/internal/apperr"
	"book-catalog/pkg/utils/mask"
	"fmt"
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
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
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
	UserStatusActive      UserStatus = ""
	UserStatusNonActivate UserStatus = "user_status_non_activate"
)

// GetAppError by UserNotActiveStatus
func (u *User) GetAppError() error {
	switch u.Data.Status {
	case UserStatusActive:
		return nil
	case UserStatusNonActivate:
		return apperr.ErrUserNotActivated
	default:
		return apperr.ErrForbidden
	}
}
