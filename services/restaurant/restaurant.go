package restaurant

import (
	"gopkg.in/guregu/null.v3"
)

// Restaurant retrieved restaurant
type Restaurant struct {
	ID            string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
}

type Table struct {
	ID            string `json:"id" db:"id"`
	RestaurantID            string `json:"restaurant_id" db:"restaurant_id"`
	Name  null.String `json:"name" db:"name"`
	NumSeatsAvailable int `json:"num_seats_available" db:"num_seats_available"`
	NumSeatsReserved int `json:"num_seats_reserved" db:"num_seats_reserved"`
	StartDate string `json:"start_date" db:"start_date"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
}
