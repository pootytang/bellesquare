package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PayoutPeriod int

const (
	FirstQuarter PayoutPeriod = iota
	SecondQuarter
	ThirdQuarter
	FourthQuarter
	FinalScore
)

// These strings match your DB's TEXT column perfectly
var payoutPeriodNames = [...]string{
	"q1",
	"q2",
	"q3",
	"q4",
	"final",
}

func (p PayoutPeriod) String() string {
	// be careful logging in here. See comments in board_status_enum --> MarshalJSON
	if p < 0 || int(p) >= len(payoutPeriodNames) {
		return "Unknown"
	}
	return payoutPeriodNames[p]
}

func (p PayoutPeriod) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p *PayoutPeriod) Scan(value interface{}) error {
	// slog.Debug("Scanning payout period from DB", "value", value)
	if value == nil {
		// slog.Warn("Received NULL value for payout period, defaulting to FirstQuarter")
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		// slog.Debug("Payout period value is []byte, converting to string")
		str = string(v)
	case string:
		// slog.Debug("Payout period value string", "value", v)
		str = v
	default:
		return fmt.Errorf("cannot scan %T into PayoutPeriod", value)
	}

	for i, name := range payoutPeriodNames {
		if name == str {
			// slog.Debug("Matched payout period string to enum", "string", str, "enum index", i)
			*p = PayoutPeriod(i)
			return nil
		}
	}
	return fmt.Errorf("invalid payout period: %s", str)
}

func (p PayoutPeriod) MarshalJSON() ([]byte, error) {
	// See comment in board_status_enum --> MarshalJSON()
	return json.Marshal(p.String())
}

func (p *PayoutPeriod) UnmarshalJSON(data []byte) error {
	// See comment in board_status_enum --> MarshalJSON()
	// slog.Debug("unmarshalling json data")
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		// slog.Error("error unmarshalling json data", "error", err)
		return err
	}

	for i, name := range payoutPeriodNames {
		if name == str {
			*p = PayoutPeriod(i)
			// slog.Debug("period name found", "name", name)
			return nil
		}
	}

	// slog.Error("period string invalid")
	return fmt.Errorf("invalid payout period string: %s", str)
}

// UnmarshalText handles the conversion when the type is used as a MAP KEY
func (p *PayoutPeriod) UnmarshalText(text []byte) error {
	// slog.Debug("UnmarshalText called")
	str := string(text)
	for i, name := range payoutPeriodNames {
		if name == str {
			*p = PayoutPeriod(i)
			// slog.Debug("period name found", "name", name)
			return nil
		}
	}

	// slog.Error("period string invalid", "name", text)
	return fmt.Errorf("invalid payout period key: %s", str)
}

func (p PayoutPeriod) DisplayName() string {
	prettyNames := [...]string{"1st Quarter", "2nd Quarter", "3rd Quarter", "4th Quarter", "Final Score"}
	if p < 0 || int(p) >= len(prettyNames) {
		return "Unknown"
	}
	return prettyNames[p]
}
