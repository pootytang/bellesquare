package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// BECAREFUL when or if logging the methods because it can cause an infinite loop
// See the comment in board_status_enum --> MarshalJSON

type SportEnum int

const (
	Unknown_Sport SportEnum = iota
	Football
	Basketball
	Baseball
	Hockey
	Soccer
	Volleyball
)

var sportNames = [...]string{"unknown", "football", "basketball", "baseball", "hockey", "soccer", "volleyball"}

// String returns the string representation for logging or internal use
func (s SportEnum) String() string {
	if s < 0 || int(s) >= len(sportNames) {
		return sportNames[Unknown_Sport]
	}
	return sportNames[s]
}

// Value tells pgx how to save this to the DB (converts int -> string)
func (s SportEnum) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan tells pgx how to read this from the DB (converts string -> int)
func (s *SportEnum) Scan(value interface{}) error {
	if value == nil {
		*s = Unknown_Sport
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot scan %T into SportEnum", value)
	}

	*s = FromString(str)
	return nil
}

// MarshalJSON ensures SvelteKit receives "football" instead of 1
func (s SportEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON ensures Go can read "football" from a JSON request
func (s *SportEnum) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*s = FromString(str)
	return nil
}

// Equals allows for case-insensitive comparison against a string
func (s SportEnum) Equals(other string) bool {
	return strings.EqualFold(s.String(), other)
}

// IsValid checks if the sport is part of our supported list
func (s SportEnum) IsValid() bool {
	// Exclude Unknown_Sport (0)
	return s > Unknown_Sport && int(s) < len(sportNames)
}

/********** HELPER FUNCTIONS **********/
func FromString(str string) SportEnum {
	cleanStr := strings.ToLower(strings.TrimSpace(str))
	for i, name := range sportNames {
		if name == cleanStr {
			return SportEnum(i)
		}
	}
	return Unknown_Sport
}
