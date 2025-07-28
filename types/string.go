package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type String struct {
	Val   string
	Valid bool
}

// NewString creates a new String from a string.
func NewString(s string) String {
	return String{Val: s, Valid: true}
}

// Scan implements the Scanner interface.
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

// Value implements the Valuer interface.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Val, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Val)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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

// IsZero returns true if the string is not valid or is an empty string.
func (s String) IsZero() bool {
	return !s.Valid || s.Val == ""
}

// String returns the string value or an empty string if not valid.
// This method implements the fmt.Stringer interface.
func (s String) String() string {
	if !s.Valid {
		return ""
	}
	return s.Val
}

// Ptr returns a pointer to the string value, or nil if not valid.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.Val
}
