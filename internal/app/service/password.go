package service

import (
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService is a service for password.
type PasswordService struct {
	cost int
}

// NewPasswordService creates a new PasswordService instance.
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost,
	}
}

func (s *PasswordService) Validate(input, hash types.Password) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	if err != nil {
		return err
	}

	return nil
}

func (s *PasswordService) Hash(password types.Password) (types.Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}

	return types.Password(hash), nil
}
