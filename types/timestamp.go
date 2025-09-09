package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp is a custom type for handling full date-time values (with timezone),
// stored in RFC3339 format. It includes a validity flag to support NULL-like
// semantics for databases and JSON.
type Timestamp struct {
	Time  time.Time // The stored timestamp value, normalized to UTC
	Valid bool
}

// Defines the standard format for timestamps (RFC3339).
const timestampFormat = time.RFC3339

// NewTimestamp creates a new valid Timestamp from a time.Time,
// normalizing to UTC and truncating to the nearest second.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{
		Time:  t.UTC().Truncate(time.Second),
		Valid: true,
	}
}

// CombineDateAndTime creates a new valid Timestamp from separate Date and Time values,
func CombineDateAndTime(d Date, t Time) Timestamp {
	date := d.Time
	tod := t.Time

	return Timestamp{
		Time: time.Date(
			date.Year(), date.Month(), date.Day(),
			tod.Hour(), tod.Minute(), tod.Second(), tod.Nanosecond(),
			date.Location(),
		),
		Valid: true,
	}
}

// Scan implements the sql.Scanner interface.
// It converts database values into a Timestamp, handling NULL, time.Time,
// []byte, and string values.
func (t *Timestamp) Scan(value any) error {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v.UTC().Truncate(time.Second)
		t.Valid = true
		return nil
	case []byte:
		return t.parseTimestampString(string(v))
	case string:
		return t.parseTimestampString(v)
	default:
		return fmt.Errorf("cannot scan %T into Timestamp", value)
	}
}

// parseTimestampString parses an RFC3339-formatted string into a Timestamp.
// If the string is empty, the Timestamp is set invalid.
func (t *Timestamp) parseTimestampString(s string) error {
	if s == "" {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	parsed, err := time.Parse(timestampFormat, s)
	if err != nil {
		return fmt.Errorf("invalid timestamp format, expected RFC3339: %w", err)
	}
	t.Time = parsed.UTC().Truncate(time.Second)
	t.Valid = true
	return nil
}

// Value implements the driver.Valuer interface.
// It converts the Timestamp into a database-compatible value (time.Time or NULL).
func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.UTC().Truncate(time.Second), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It converts the Timestamp into a JSON string in RFC3339 format, or null if invalid.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UTC().Truncate(time.Second).Format(timestampFormat))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It parses a JSON string into a Timestamp, handling null and empty strings.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	// Remove surrounding quotes if present
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return t.parseTimestampString(str)
}

// IsZero reports whether the Timestamp is invalid or represents the zero time.
func (t Timestamp) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// String returns the Timestamp formatted in RFC3339, or an empty string if invalid.
// Implements the fmt.Stringer interface.
func (t Timestamp) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timestampFormat)
}
