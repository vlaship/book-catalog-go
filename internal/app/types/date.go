package types

import (
	"fmt"
	"time"
)

const (
	dayFormat = "2006-01-02"
)

// DateDay is a custom YYYY-MM-DD type to handle JSON serialization and deserialization.
type DateDay struct {
	time.Time
}

// MarshalJSON to customize the JSON encoding for DateDay.
func (ct *DateDay) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("%q", ct.Time.Format(dayFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON to customize the JSON decoding for DateDay.
func (ct *DateDay) UnmarshalJSON(data []byte) error {
	// Remove the surrounding double quotes before parsing the time string.
	timeString := string(data[1:11])
	parsedTime, err := time.Parse(dayFormat, timeString)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}
