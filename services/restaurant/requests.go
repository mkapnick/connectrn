package restaurant

// RestaurantCreateRequest create a restaurant
type RestaurantCreateRequest struct {
	Name string `validate:"required" json:"name"`
}
