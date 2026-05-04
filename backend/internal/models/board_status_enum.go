package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type BoardStatus int

const (
	BoardOpen BoardStatus = iota
	BoardLocked
	BoardCompleted
)

var boardStatusNames = [...]string{"open", "locked", "completed"}

func (s BoardStatus) String() string {
	// See comment in MarshalJSON below
	// slog calls this same BoardStatus which again calls MarshalJSON
	// slog.Debug("Converting BoardStatus to string", "status", s)
	if s < 0 || int(s) >= len(boardStatusNames) {
		return "unknown"
	}
	return boardStatusNames[s]
}

func (s BoardStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *BoardStatus) Scan(value interface{}) error {
	// slog.Debug("Scanning BoardStatus from database value", "value", value)
	if value == nil {
		// slog.Debug("Received NULL value for BoardStatus, defaulting to BoardOpen")
		*s = BoardOpen
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		// slog.Debug("BoardStatus value is []byte, converting to string", "value", v)
		str = string(v)
	case string:
		// slog.Debug("BoardStatus value is string", "value", v)
		str = v
	default:
		return fmt.Errorf("cannot scan %T into BoardStatus", value)
	}

	for i, name := range boardStatusNames {
		if name == str {
			// slog.Debug("Matched BoardStatus string to enum", "string", str, "enum", i)
			*s = BoardStatus(i)
			return nil
		}
	}
	return fmt.Errorf("invalid board status: %s", str)
}

// JSON Marshalling for SvelteKit
func (s BoardStatus) MarshalJSON() ([]byte, error) {
	// This can cause an infinite loop because slog itself is logging using json and calls this method again and again
	// Need to either remove the slog line or force the string so MarshalJSON is not called
	// slog.Debug("Marshalling BoardStatus to JSON", "status", string(boardStatusNames[s]))
	return json.Marshal(s.String())
}

func (s *BoardStatus) UnmarshalJSON(data []byte) error {
	// See comment in MarshalJSON - this wasn't failing but to be safe, commenting the slog entry out
	// slog.Debug("Unmarshalling BoardStatus from JSON", "data", data)
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	for i, name := range boardStatusNames {
		if name == str {
			*s = BoardStatus(i)
			return nil
		}
	}
	return fmt.Errorf("invalid board status: %s", str)
}
