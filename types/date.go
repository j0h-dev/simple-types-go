package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	Time  time.Time
	Valid bool
}

const dateFormat = "2006-01-02"

func NewDate(t time.Time) Date {
	return Date{Time: t.Truncate(24 * time.Hour), Valid: true}
}

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

func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time.Format(dateFormat), nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	str := fmt.Sprintf(`"%s"`, d.Time.Format(dateFormat))
	return []byte(str), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		d.Time, d.Valid = time.Time{}, false
		return nil
	}

	// Remove quotes if present
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return d.parseDateString(str)
}

func (d Date) IsZero() bool {
	return !d.Valid || d.Time.IsZero()
}

func (d Date) String() string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format(dateFormat)
}
