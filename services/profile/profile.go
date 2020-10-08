package profile

import (
	"gopkg.in/guregu/null.v3"
)

// Profile is a retrieved profile
type Profile struct {
	ID        string      `json:"id" db:"id"`
	AccountID string      `json:"account_id" db:"account_id"`
	Name      null.String `json:"name" db:"name"`
	CreatedAt string      `json:"created_at" db:"created_at"`
	UpdatedAt string      `json:"updated_at" db:"updated_at"`

	// this is a db column on the `account` table
	Email string `json:"email" db:"email"`
}
