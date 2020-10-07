package account

import (
	"gopkg.in/guregu/null.v3"
)

// Account is a retrieved and account
type Account struct {
	ID         string      `json:"id" db:"id"`
	RestaurantID  null.String `json:"company_id" db:"company_id"`
	Email      string      `json:"email" db:"email"`
	Password   string      `json:"-" db:"password"`
	CreatedAt  string      `json:"created_at" db:"created_at"`
	UpdatedAt  string      `json:"updated_at" db:"updated_at"`

	// this is ONLY used to attach the `profile_id` to the account object.
	// this field does NOT exist as a db field
	ProfileID string `json:"profile_id"`
}

// Restaurant retrieves a company subdomain
type Restaurant struct {
	ID        string `json:"id" db:"id"`
	Subdomain string `json:"subdomain" db:"subdomain"`
}
