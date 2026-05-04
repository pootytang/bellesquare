package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `form:"email" json:"email" db:"email"`
	UserName  string    `form:"username" json:"user_name" db:"username"`
	FirstName string    `form:"first_name" json:"first_name" db:"first_name"`
	LastName  string    `form:"last_name" json:"last_name" db:"last_name"`
	IsGuest   bool      `json:"is_guest" db:"is_guest"`
	// db:"-" tells pgx to ignore this field for SQL queries
	Password       string    `form:"password" db:"-"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	VenmoHandle    string    `form:"venmo_handle" json:"venmo_handle" db:"venmo_handle"`
	ZelleHandle    string    `form:"zelle_handle" json:"zelle_handle" db:"zelle_handle"`
	PreferredColor string    `form:"preferred_color" json:"preferred_color" db:"preferred_color"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
