package request

import (
	"fmt"
	"github.com/vlaship/book-catalog-go/pkg/utils/mask"
)

// UserData is a request model for user info
type UserData struct {
	FirstName string `json:"firstname" validate:"required,min=2,max=64" example:"John"`
	LastName  string `json:"lastname" validate:"required,min=2,max=64" example:"Doe"`
	Email     string `json:"email" validate:"required,email" example:"email@email.com"`
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
