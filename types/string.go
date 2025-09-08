package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// String is a custom type for handling nullable strings.
// It wraps a string value and a validity flag, similar to sql.NullString,
// but with extra helpers for JSON and convenience.
type String struct {
	Val   string
	Valid bool
}

// Creates a new valid String from a raw string.
func NewString(s string) String {
	return String{Val: s, Valid: true}
}

// Scan implements the sql.Scanner interface.
// It converts database values into a String, supporting NULL, string, and []byte.
func (s *String) Scan(value any) error {
	if value == nil {
		s.Val, s.Valid = "", false
		return nil
	}

	switch v := value.(type) {
	case string:
		s.Val = v
		s.Valid = true
		return nil
	case []byte:
		s.Val = string(v)
		s.Valid = true
		return nil
	default:
		return fmt.Errorf("cannot scan %T into String", value)
	}
}

// Value implements the driver.Valuer interface.
// It returns the string value for database storage, or nil if invalid.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Val, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It encodes the string as a JSON string, or null if invalid.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Val)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It decodes JSON input into the String type, handling "null" as invalid.
func (s *String) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.Val, s.Valid = "", false
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid string format: %w", err)
	}
	s.Val = str
	s.Valid = true
	return nil
}

// IsZero returns true if the String is invalid or contains an empty string.
// Useful for omitempty behavior in JSON or zero-value checks.
func (s String) IsZero() bool {
	return !s.Valid || s.Val == ""
}

// String returns the underlying string value, or an empty string if invalid.
// Implements the fmt.Stringer interface.
func (s String) String() string {
	if !s.Valid {
		return ""
	}
	return s.Val
}

// Ptr returns a pointer to the underlying string value.
// Returns nil if the String is invalid. Useful for APIs expecting *string.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.Val
}
