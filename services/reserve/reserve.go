package reserve

// Reserve retrieved reserve date
type Reserve struct {
	ID            string `json:"id" db:"id"`
	GolfCourseID  string `json:"golf_course_id" db:"golf_course_id"`
	StartDate     string `json:"start_date" db:"start_date"`
	EndDate       string `json:"end_date" db:"end_date"`
	Reason        string `json:"reason" db:"reason"`
	Name          string `json:"name" db:"name"`
	CreatedBy     string `json:"created_by" db:"created_by"`
	LastUpdatedBy string `json:"last_updated_by" db:"last_updated_by"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
}
