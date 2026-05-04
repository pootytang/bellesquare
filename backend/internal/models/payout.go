package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Payout struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	BoardID      uuid.UUID    `json:"board_id" db:"board_id"`
	PeriodName   PayoutPeriod `json:"period_name" db:"period_name"` // Typed Enum
	Amount       float64      `json:"amount" db:"amount"`           // NUMERIC(10,2)
	WinnerUserID *uuid.UUID   `json:"winner_user_id" db:"winner_user_id"`
	WinnerName   *string      `json:"winner_name"`
	AwardedAt    *time.Time   `json:"awarded_at" db:"awarded_at"`
}
