package request

import (
	"fmt"
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"github.com/vlaship/go-mask"
)

// Signin request
type Signin struct {
	Username types.Username `json:"username" validate:"required,email"`
	Password types.Password `json:"password" validate:"required,min=8,max=64"`
}

// String
func (s *Signin) String() string {
	return fmt.Sprintf("Signin{Username: %s, Password: %s}", mask.String(s.Username.String()), mask.String(s.Password.String()))
}

// Signup request
type Signup struct {
	Username  types.Username `json:"username" validate:"required,email"`
	Password  types.Password `json:"password" validate:"required,min=8,max=64"`
	Firstname string         `json:"firstname" validate:"required,min=2,max=64"`
	Lastname  string         `json:"lastname"`
}

// String
func (s *Signup) String() string {
	return fmt.Sprintf("Signup{Username: %s, Password: %s}", mask.String(s.Username.String()), mask.String(s.Password.String()))
}

// Activation request
type Activation struct {
	OTP types.Token `json:"otp" validate:"required,min=64,max=64"`
}

// String
func (a *Activation) String() string {
	return fmt.Sprintf("Activation{Token: %s}", mask.String(string(a.OTP)))
}

// ResendActivation request
type ResendActivation struct {
	Username types.Username `json:"username" validate:"required,email"`
}

// String
func (r *ResendActivation) String() string {
	return fmt.Sprintf("Resend{Username: %s}", mask.String(r.Username.String()))
}

// ResetPassword request
type ResetPassword struct {
	Username types.Username `json:"username" validate:"required,email"`
}

// String
func (r *ResetPassword) String() string {
	return fmt.Sprintf("ResetPassword{Username: %s}", mask.String(r.Username.String()))
}

// ReplacePassword request
type ReplacePassword struct {
	OTP         types.Token    `json:"otp" validate:"required,min=64,max=64"`
	NewPassword types.Password `json:"new_password" validate:"required,min=8,max=64"`
}

// String
func (r *ReplacePassword) String() string {
	return fmt.Sprintf(
		"ReplacePassword{OTP: %s, NewPassword: %s}",
		mask.String(string(r.OTP)),
		mask.String(r.NewPassword.String()),
	)
}

// ChangePassword request
type ChangePassword struct {
	CurrentPassword types.Password `json:"current_password" validate:"required,min=8,max=64"`
	NewPassword     types.Password `json:"new_password" validate:"required,min=8,max=64"`
}

// String
func (u *ChangePassword) String() string {
	return fmt.Sprintf(
		"ChangePassword{CurrentPassword: %s, NewPassword: %s}",
		mask.String(u.CurrentPassword.String()),
		mask.String(u.NewPassword.String()),
	)
}
