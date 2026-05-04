package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Team struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CreatorID      uuid.UUID `json:"creator_id" db:"creator_id"`
	City           string    `json:"city" db:"city"`     // New: "Detroit"
	Mascot         string    `json:"mascot" db:"mascot"` // Renamed: "Lions"
	Sport          SportEnum `json:"sport" db:"sport"`
	TeamLogoURL    *string   `json:"team_logo_url" db:"team_logo_url"`
	PrimaryColor   string    `json:"primary_color" db:"primary_color"`
	SecondaryColor string    `json:"secondary_color" db:"secondary_color"`
	TertiaryColor  *string   `json:"tertiary_color" db:"tertiary_color"` // Pointer for optional
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Helper method for the UI
func (t Team) FullName() string {
	return t.City + " " + t.Mascot
}

func (t Team) String() string {
	// Example output: "Detroit Lions (football)"
	return fmt.Sprintf("%s %s (%s)", t.City, t.Mascot, t.Sport.String())
}

type TeamCounts struct {
	Football   int `json:"football"`
	Basketball int `json:"basketball"`
	Baseball   int `json:"baseball"`
	Soccer     int `json:"soccer"`
	Hockey     int `json:"hockey"`
}
