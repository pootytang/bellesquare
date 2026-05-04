package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Board struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	Title          string      `json:"title" db:"title"`
	Sport          SportEnum   `json:"sport" db:"sport"`
	Status         BoardStatus `json:"status" db:"status"`
	PricePerSquare float64     `json:"price_per_square" db:"price_per_square"`

	// Postgres arrays map directly to Go slices
	HomeAxisNumbers []int32 `json:"home_axis_numbers" db:"home_axis_numbers"`
	AwayAxisNumbers []int32 `json:"away_axis_numbers" db:"away_axis_numbers"`

	CreatorID  uuid.UUID `json:"creator_id" db:"creator_id"`
	HomeTeamID uuid.UUID `json:"home_team_id" db:"home_team_id"`
	AwayTeamID uuid.UUID `json:"away_team_id" db:"away_team_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (b Board) String() string {
	return fmt.Sprintf("Board: %s [%s] - Price: $%.2f", b.Title, b.Status.String(), b.PricePerSquare)
}

// GenerateRandomAxes populates the 0-9 slices for the board
func (b *Board) GenerateRandomAxes() {
	home := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	away := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// rand.Seed(time.Now().UnixNano()) // Not needed since Go 1.20, the default Source is already seeded with a good seed based on time.
	rand.Shuffle(len(home), func(i, j int) { home[i], home[j] = home[j], home[i] })
	rand.Shuffle(len(away), func(i, j int) { away[i], away[j] = away[j], away[i] })

	b.HomeAxisNumbers = home
	b.AwayAxisNumbers = away
}

type BoardPlayer struct {
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Color       string    `json:"color"`
	SquareCount int       `json:"square_count"`
	IsPaid      bool      `json:"is_paid"` // Only populated for Creator
}
