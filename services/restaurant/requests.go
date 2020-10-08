package restaurant

// RestaurantCreateRequest create a restaurant
type RestaurantCreateRequest struct {
	Name string `validate:"required" json:"name"`
}

// TableCreateRequest create a table at a restaurant
type TableCreateRequest struct {
	RestaurantID string `validate:"required" json:"restaurant_id"`
	Name string `validate:"required" json:"name"`
	NumSeatsAvailable int `validate:"required,lte=4" json:"num_seats_available"`
	NumSeatsReserved int `json:"num_seats_reserved"`
	StartDate string `validate:"required" json:"start_date"`
}
