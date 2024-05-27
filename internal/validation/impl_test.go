package validation

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimalMinValidation(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-min=10.0"`
	}{})
	assert.Error(t, err)
}

func TestDecimalMinValidationSuccess(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-min=10.0"`
	}{Decimal: decimal.NewFromFloat(10.0)})
	assert.NoError(t, err)
}

func TestDecimalMaxValidation(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-max=10.0"`
	}{Decimal: decimal.NewFromFloat(11.0)})
	assert.Error(t, err)
}

func TestDecimalMaxValidationSuccess(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-max=10.0"`
	}{Decimal: decimal.NewFromFloat(10.0)})
	assert.NoError(t, err)
}

func TestDecimalPositiveMinusValidation(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-positive"`
	}{Decimal: decimal.NewFromFloat(-1.0)})
	assert.Error(t, err)
}

func TestDecimalPositiveZeroValidation(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-positive"`
	}{Decimal: decimal.NewFromFloat(0.0)})
	assert.Error(t, err)
}

func TestDecimalPositiveValidationSuccess(t *testing.T) {
	validator := New()
	err := validator.Struct(struct {
		Decimal decimal.Decimal `validate:"dec-positive"`
	}{Decimal: decimal.NewFromFloat(1.0)})
	assert.NoError(t, err)
}
