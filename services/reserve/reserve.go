package reserve

// UserReservation user reservation
type UserReservation struct {
	ID           string `json:"id" db:"id"`
	RestaurantID string `json:"restaurant_id" db:"restaurant_id"`
	TableID      string `json:"table_id" db:"table_id"`
	ProfileID    string `json:"profile_id" db:"profile_id"`
	NumSeats     int `json:"num_seats" db:"num_seats"`
	StartDate    string `json:"start_date" db:"start_date"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

// UserReservationCanceled canceled user reservation
type UserReservationCanceled struct {
	ID           string `json:"id" db:"id"`
	RestaurantID string `json:"restaurant_id" db:"restaurant_id"`
	TableID      string `json:"table_id" db:"table_id"`
	ProfileID    string `json:"profile_id" db:"profile_id"`
	NumSeats     string `json:"num_seats" db:"num_seats"`
	StartDate    string `json:"start_date" db:"start_date"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

// Table representing a restaurant table
type Table struct {
	ID                string `json:"id" db:"id"`
	RestaurantID      string `json:"restaurant_id" db:"restaurant_id"`
	Name              string `json:"name" db:"name"`
	NumSeatsAvailable int    `json:"num_seats_available" db:"num_seats_available"`
	NumSeatsReserved  int    `json:"num_seats_reserved" db:"num_seats_reserved"`
	StartDate         string `json:"start_date" db:"start_date"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	UpdatedAt         string `json:"updated_at" db:"updated_at"`
}
