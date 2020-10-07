package profile

import (
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
)

// Profile is a retrieved profile
type Profile struct {
	ID               string         `json:"id" db:"id"`
	AccountID        string         `json:"account_id" db:"account_id"`
	FirstName        string         `json:"first_name" db:"first_name"`
	LastName         string         `json:"last_name" db:"last_name"`
	PhoneNumber      string         `json:"phone_number" db:"phone_number"`
	Handicap         null.String    `json:"handicap" db:"handicap"`
	NumRoundsPerYear null.String    `json:"num_rounds_per_year" db:"num_rounds_per_year"`
	Preferences      types.JSONText `json:"preferences" db:"preferences"`
	InviteCode       null.String    `json:"invite_code" db:"invite_code"`
	InvitedBy        null.String    `json:"invited_by" db:"invited_by"`
	InvitedAt        null.String    `json:"invited_at" db:"invited_at"`
	CreatedAt        string         `json:"created_at" db:"created_at"`
	UpdatedAt        string         `json:"updated_at" db:"updated_at"`

	// this is a db column on the `account` table
	Email string `json:"email" db:"email"`
}
