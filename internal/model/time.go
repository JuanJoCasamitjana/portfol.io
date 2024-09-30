package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type RFC3339NanoTime struct {
	time.Time
}

// Scan is the method used to convert a database value into a Go value
func (t *RFC3339NanoTime) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan time: %v", value)
	}
	parsedTime, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			return err
		}
	}
	t.Time = parsedTime
	return nil
}

// Value converts the Go type to a database value
func (t RFC3339NanoTime) Value() (driver.Value, error) {
	return t.Time.Format(time.RFC3339Nano), nil
}
