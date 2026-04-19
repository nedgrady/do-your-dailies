package utctime

import (
	"encoding/json"
	"fmt"
	"time"
)

// Time normalizes all values to UTC for API boundaries.
type Time struct {
	time.Time
}

func (value *Time) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return fmt.Errorf("parse date-time: %w", err)
	}

	value.Time = parsed.UTC()
	return nil
}

func (value Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.UTC().Format(time.RFC3339Nano))
}
