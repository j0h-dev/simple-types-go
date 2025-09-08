package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Date is a custom type for representing dates (without time-of-day).
type Date struct {
	Time  time.Time
	Valid bool
}

// Defines the standard format for dates (YYYY-MM-DD).
const dateFormat = "2006-01-02"

// NewDate creates a new valid Date, truncating the time to midnight.
func NewDate(t time.Time) Date {
	return Date{Time: t.Truncate(24 * time.Hour), Valid: true}
}

// Scan implements the sql.Scanner interface.
// It converts a database value into a Date, handling NULL, time.Time, []byte, and string inputs.
func (d *Date) Scan(value any) error {
	if value == nil {
		d.Time, d.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v.Truncate(24 * time.Hour)
		d.Valid = true
		return nil
	case []byte:
		return d.parseDateString(string(v))
	case string:
		return d.parseDateString(v)
	default:
		return fmt.Errorf("cannot scan %T into Date", value)
	}
}

// Parses a string in YYYY-MM-DD format into a Date.
// If the string is empty, the Date is marked invalid.
func (d *Date) parseDateString(s string) error {
	if s == "" {
		d.Time, d.Valid = time.Time{}, false
		return nil
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	d.Valid = true
	return nil
}

// Value implements the driver.Valuer interface.
// It converts the Date into a database-compatible value (string or NULL).
func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time.Format(dateFormat), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It converts the Date into a JSON string (or null if invalid).
func (d Date) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	str := fmt.Sprintf(`"%s"`, d.Time.Format(dateFormat))
	return []byte(str), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It parses a JSON string into a Date, handling null and empty strings.
func (d *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		d.Time, d.Valid = time.Time{}, false
		return nil
	}

	// Remove surrounding quotes if present
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return d.parseDateString(str)
}

// IsZero reports whether the Date is invalid or represents the zero time.
func (d Date) IsZero() bool {
	return !d.Valid || d.Time.IsZero()
}

// String returns the Date formatted as YYYY-MM-DD, or an empty string if invalid.
func (d Date) String() string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format(dateFormat)
}
