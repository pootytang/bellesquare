package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Square struct {
	ID       uuid.UUID `json:"id" db:"id"`
	BoardID  uuid.UUID `json:"board_id" db:"board_id"`
	RowIndex int       `json:"row_index" db:"row_index"`
	ColIndex int       `json:"col_index" db:"col_index"`

	UserID    *uuid.UUID `json:"user_id" db:"user_id"`
	UserColor string     `json:"user_color" db:"user_color"`
	// New field for the UI
	UserInitials string `json:"user_initials" db:"user_initials"`

	IsPaid        bool          `json:"is_paid" db:"is_paid"`
	PaymentStatus PaymentStatus `json:"payment_status" db:"payment_status"`
	ClaimedAt     *time.Time    `json:"claimed_at" db:"claimed_at"`
}

// Used by the websocket
type SquareUpdate struct {
	Row           int           `json:"row"`
	Col           int           `json:"col"`
	UserID        uuid.UUID     `json:"user_id"`
	UserColor     string        `json:"user_color"`
	UserInitials  string        `json:"user_initials"`
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	IsPaid        bool          `json:"is_paid"`
}

type SquareCoords struct {
	Row int `json:"row"`
	Col int `json:"col"`
}
