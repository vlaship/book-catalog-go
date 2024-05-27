package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"reflect"
)

// ValidatorImpl is a validator for requests
type ValidatorImpl struct {
	valid *validator.Validate
}

// Struct validates a struct
func (v *ValidatorImpl) Struct(a any) error {
	return v.valid.Struct(a)
}

func New() Validator {
	v := &ValidatorImpl{
		valid: validator.New(),
	}
	_ = v.decimalMin()
	_ = v.decimalMax()
	_ = v.decimalPositive()
	return v
}

func (v *ValidatorImpl) decimalMin() error {
	v.valid.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.String()
		}
		return nil
	}, decimal.Decimal{})
	return v.valid.RegisterValidation("dec-min", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.GreaterThanOrEqual(baseValue)
	})
}

func (v *ValidatorImpl) decimalMax() error {
	v.valid.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.String()
		}
		return nil
	}, decimal.Decimal{})
	return v.valid.RegisterValidation("dec-max", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.LessThanOrEqual(baseValue)
	})
}

func (v *ValidatorImpl) decimalPositive() error {
	v.valid.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.String()
		}
		return nil
	}, decimal.Decimal{})
	return v.valid.RegisterValidation("dec-positive", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		return value.GreaterThan(decimal.Zero)
	})
}
