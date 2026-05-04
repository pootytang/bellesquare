package models

import "github.com/gofrs/uuid/v5"

//TODO: Merge CreateTeamRequest and TeamDTO - TeamRequestDTO. Difference is SportEnum. I like using the enum better
/********** TEAM DTOs **********/
type CreateTeamRequest struct {
	CreatorID      uuid.UUID `json:"creator_id"`
	City           string    `json:"city"`
	Mascot         string    `json:"mascot"`
	Sport          SportEnum `json:"sport"`
	TeamLogoURL    *string   `json:"team_logo_url"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor string    `json:"secondary_color"`
}

type TeamDTO struct {
	ID             string  `json:"id"`
	CreatorID      string  `json:"creator_id"`
	FullName       string  `json:"full_name"`
	City           string  `json:"city"`
	Mascot         string  `json:"mascot"`
	Sport          string  `json:"sport"`
	TeamLogoURL    *string `json:"team_logo_url"`
	PrimaryColor   string  `json:"primary_color"`
	SecondaryColor string  `json:"secondary_color"`
}

type TeamSummary struct {
	ID             uuid.UUID `json:"id"`
	City           string    `json:"city"`
	FullName       string    `json:"full_name"`
	Mascot         string    `json:"mascot"`
	TeamLogoURL    *string   `json:"team_logo_url"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor string    `json:"secondary_color"`
}

// MapToDTO converts a database Team model to a TeamDTO
func (t Team) MapToDTO() TeamDTO {
	return TeamDTO{
		ID:             t.ID.String(),
		CreatorID:      t.CreatorID.String(),
		FullName:       t.FullName(), // Uses your existing helper method
		City:           t.City,
		Mascot:         t.Mascot,
		Sport:          t.Sport.String(),
		TeamLogoURL:    t.TeamLogoURL,
		PrimaryColor:   t.PrimaryColor,
		SecondaryColor: t.SecondaryColor,
	}
}
