package restaurant

// RestaurantsRequest request to get golf course restaurants
type RestaurantsRequest struct {
	GolfCourseID string `validate:"required" json:"golf_course_id"`
	StartDate    string `validate:"required" json:"start_date"`
	EndDate      string `validate:"required" json:"end_date"`
}
