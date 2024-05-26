package validation

// Validator is a validator interface.
//
//go:generate mockgen -destination=../../test/mock/validation/mock-validator.go -package=mock . Validator
type Validator interface {
	Struct(any) error
}
