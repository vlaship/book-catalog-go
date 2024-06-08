package types

import (
	"fmt"
	"github.com/vlaship/book-catalog-go/pkg/utils/mask"
	"strconv"
)

// UserID is a custom type for a user ID
type UserID int64

// String to string
func (u *UserID) String() string {
	return fmt.Sprintf("%d", *u)
}

// NewUserID creates a new UserID
func NewUserID(id string) (UserID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return UserID(i), nil
}

// ID is a custom type for an ID
type ID int64

// String to string
func (i *ID) String() string {
	return fmt.Sprintf("%d", *i)
}

// NewID creates a new ID
func NewID(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return ID(i), nil
}

// Token is a custom type for a token
type Token string

// Username is a custom type for a username
type Username string

// String to string
func (u *Username) String() string {
	return mask.String(string(*u))
}

// Password is a custom type for a password
type Password string

// String to string
func (p *Password) String() string {
	return mask.String(string(*p))
}
