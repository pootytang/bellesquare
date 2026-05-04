package models

import (
	"github.com/gofrs/uuid/v5"
)

/********** USER DTOs **********/
// CreateUserRequest defines what we expect from the frontend, The API contract
// Note: This is separate from the User model because we want to catch the raw password here before hashing it
type RegisterUserRequest struct {
	Email          string `form:"email"`
	Username       string `form:"user_name"`
	FirstName      string `form:"first_name"`
	LastName       string `form:"last_name"`
	Password       string `form:"password"`
	IsGuest        bool   `form:"is_guest"`
	VenmoHandle    string `form:"venmo_handle"`
	ZelleHandle    string `form:"zelle_handle"`
	PreferredColor string `form:"preferred_color"`
}

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UpdatePasswordRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	OldPassword string    `json:"old_password"`
	NewPassword string    `json:"new_password"`
}

type UpdateRequest struct {
	Username       string `json:"user_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	VenmoHandle    string `json:"venmo_handle"`
	ZelleHandle    string `json:"zelle_handle"`
	PreferredColor string `json:"preferred_color"`
	IsGuest        string `json:"is_guest"`
}

type UpgradeRequest struct {
	UserID       uuid.UUID `json:"-"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"-"` // This will be calculated in the handler, not sent from the frontend
}

/********** SQUARE DTOs **********/
type ClaimSquareRequest struct {
	UserID     uuid.UUID      `json:"user_id"`
	UserColor  string         `json:"user_color"` // Match the 'color' field from Svelte
	Selections []SquareCoords `json:"selections"`
}

/********** PAYOUTS DTOs **********/
type UpdatePayoutsReq struct {
	// Go will use UnmarshalText to turn keys like "q1" into PayoutPeriod(0)
	Amounts map[PayoutPeriod]float64 `json:"amounts"`
}

type PayoutSaveRequest struct {
	Payouts []Payout `json:"payouts"`
}
