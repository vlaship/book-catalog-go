package types

import (
	"strings"

	"github.com/shopspring/decimal"
)

const (
	decimalPlaces = 2
)

// Decimal is a custom type for money
type Decimal struct {
	decimal.Decimal
}

// String to string
func (d *Decimal) String() string {
	return d.StringFixedBank(2)
}

// MarshalJSON to customize the JSON encoding for DateDay.
func (d *Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.StringFixedBank(decimalPlaces)), nil
}

// UnmarshalJSON to customize the JSON decoding for DateDay.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	s := strings.ReplaceAll(string(data), "\"", "")
	v, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	d.Decimal = v
	return nil
}

// PositiveDecimal is a custom type for positive decimal
type PositiveDecimal struct {
	Value decimal.Decimal // todo: `validate:"dec-positive"`
}

// String to string
func (d *PositiveDecimal) String() string {
	return d.Value.StringFixedBank(decimalPlaces)
}

// UnmarshalJSON to customize the JSON decoding for DateDay.
func (d *PositiveDecimal) UnmarshalJSON(data []byte) error {
	s := strings.ReplaceAll(string(data), "\"", "")
	v, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	d.Value = v
	return nil
}
