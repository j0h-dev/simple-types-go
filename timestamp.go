package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Timestamp struct {
	Time  time.Time
	Valid bool
}

const timestampFormat = time.RFC3339

func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{
		Time:  t.UTC().Truncate(time.Second),
		Valid: true,
	}
}

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

func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.UTC().Truncate(time.Second), nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UTC().Truncate(time.Second).Format(timestampFormat))
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return t.parseTimestampString(str)
}

func (t Timestamp) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

func (t Timestamp) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timestampFormat)
}
