package reserve

// ReserveRequest request
type ReserveRequest struct {
	TableID string `validate:"required" json:"table_id"`
	NumSeats    int `validate:"required" json:"num_seats"`
}
