package reserve

// ReserveRequest request to create a reservation
type ReserveRequest struct {
	RestaurantID string `validate:"required" json:"restaurant_id"`
	TableID string `validate:"required" json:"table_id"`
	NumSeatsReserved    int `validate:"required" json:"num_seats_reserved"`
	// tacked on to the request from the backend. Not included from client.
	ProfileID string `json:"profile_id"`
}

// CancelReserveRequest request to cancel a reservation
type CancelReserveRequest struct {
	RestaurantID string `validate:"required" json:"restaurant_id"`
	TableID string `validate:"required" json:"table_id"`
	UserReservationID string `validate:"required" json:"user_reservation_id"`
	// tacked on to the request from the backend. Not included from client.
	ProfileID string `json:"profile_id"`
}
