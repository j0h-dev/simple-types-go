package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	Time  time.Time
	Valid bool
}

const timeFormat = "15:04"

func NewTime(t time.Time) Time {
	h, m, _ := t.Clock()
	return Time{
		Time:  time.Date(1, 1, 1, h, m, 0, 0, time.UTC),
		Valid: true,
	}
}

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

func (t *Time) parseTimeString(s string) error {
	if s == "" {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

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

func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Format(timeFormat), nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	str := fmt.Sprintf(`"%s"`, t.Time.Format(timeFormat))
	return []byte(str), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return t.parseTimeString(str)
}

func (t Time) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

func (t Time) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timeFormat)
}
