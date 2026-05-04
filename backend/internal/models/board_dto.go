package models

import (
	"github.com/gofrs/uuid/v5"
)

type CreateBoardRequest struct {
	Title          string                   `json:"title"`
	Sport          SportEnum                `json:"sport"`
	PricePerSquare float64                  `json:"price_per_square"`
	HomeTeamID     uuid.UUID                `json:"home_team_id"`
	AwayTeamID     uuid.UUID                `json:"away_team_id"`
	Payouts        map[PayoutPeriod]float64 `json:"payouts"` // Map matching the PayoutPeriod Enum
}

// BoardWithSquares is a helper struct for the API response
type BoardWithSquares struct {
	Board   BoardWithTeams `json:"board"`
	Squares [10][10]Square `json:"squares"`
}

type BoardScores struct {
	Home     int                  `json:"home"`
	Away     int                  `json:"away"`
	Quarters []UpdateScoreRequest `json:"quarters"`
}

type SquareOwnerDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Color     string `json:"color"`
}

type BoardDetailsResponse struct {
	// We nest the team data here so SvelteKit gets it in one go
	Board        BoardWithTeams                `json:"board"`
	Squares      [10][10]Square                `json:"squares"`
	Payouts      []Payout                      `json:"payouts"`
	SquareOwners map[string]SquareOwnerDetails `json:"square_owners"`
	Scores       BoardScores                   `json:"scores"`
	CreatorVenmo string                        `json:"creator_venmo"`
	CreatorZelle string                        `json:"creator_zelle"`
}

// This struct acts as the "Hydrated" Board
type BoardWithTeams struct {
	// Embed the original Board fields
	*Board

	// Add the specific Team info the frontend needs
	HomeTeam TeamSummary `json:"home_team"`
	AwayTeam TeamSummary `json:"away_team"`
}

// Used by board_repo.GetBoardByUserID. This is to show the creator the boards they have available
type BoardSummary struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	Status         BoardStatus `json:"status"`
	PricePerSquare float64     `json:"price_per_square"`
	// Nest the team summaries to provide color/logo data
	HomeTeam TeamSummary `json:"home_team"`
	AwayTeam TeamSummary `json:"away_team"`
}

type BoardStatusRequest struct {
	// UnmarshalJSON will validate this automatically
	Status BoardStatus `json:"status"`
}

type UpdateScoreRequest struct {
	Quarter   int  `json:"quarter"`
	HomeScore int  `json:"home_score"`
	AwayScore int  `json:"away_score"`
	IsLocked  bool `json:"is_locked"`
}
