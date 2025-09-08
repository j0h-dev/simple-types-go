package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time is a custom type for representing only the time of day (HH:MM),
// without any associated date. It includes a validity flag to support
// NULL-like semantics for database and JSON operations.
type Time struct {
	Time  time.Time // The stored time-of-day (date is always set to year 1, month 1, day 1, UTC)
	Valid bool
}

// Defines the layout for parsing/formatting times (24-hour HH:MM).
const timeFormat = "15:04"

// NewTime creates a new valid Time from a time.Time,
// stripping away the date and seconds while keeping only HH:MM.
func NewTime(t time.Time) Time {
	h, m, _ := t.Clock()
	return Time{
		Time:  time.Date(1, 1, 1, h, m, 0, 0, time.UTC),
		Valid: true,
	}
}

// Scan implements the sql.Scanner interface.
// It converts database values into a Time, handling NULL, time.Time, []byte, and string values.
func (t *Time) Scan(value any) error {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		h, m, _ := v.Clock()
		t.Time = time.Date(1, 1, 1, h, m, 0, 0, time.UTC)
		t.Valid = true
		return nil
	case []byte:
		return t.parseTimeString(string(v))
	case string:
		return t.parseTimeString(v)
	default:
		return fmt.Errorf("cannot scan %T into Time", value)
	}
}

// parseTimeString parses a string in HH:MM format into a Time.
// If the string is empty, the Time is set invalid.
// If longer than 5 characters, only the first 5 are considered.
func (t *Time) parseTimeString(s string) error {
	if s == "" {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	// Trim to HH:MM if input includes seconds or other trailing characters
	if len(s) > 5 {
		s = s[:5]
	}

	parsed, err := time.Parse(timeFormat, s)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:MM: %w", err)
	}
	t.Time = time.Date(1, 1, 1, parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
	t.Valid = true
	return nil
}

// Value implements the driver.Valuer interface.
// It converts the Time into a database-compatible value (string or NULL).
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Format(timeFormat), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It converts the Time into a JSON string ("HH:MM") or null if invalid.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	str := fmt.Sprintf(`"%s"`, t.Time.Format(timeFormat))
	return []byte(str), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It parses a JSON string into a Time, handling null and empty strings.
func (t *Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	// Remove surrounding quotes if present
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return t.parseTimeString(str)
}

// IsZero reports whether the Time is invalid or represents the zero value.
func (t Time) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// String returns the Time formatted as "HH:MM", or an empty string if invalid.
// Implements the fmt.Stringer interface.
func (t Time) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timeFormat)
}
