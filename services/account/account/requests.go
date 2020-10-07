package account

// AccountCredentials is used when requesting in some nature
// the access to an account.
type AccountCredentials struct {
	Email     string `validate:"required" json:"email"`
	Password  string `validate:"required,gte=4" json:"password"`
	CompanyID string `json:"company_id"`
	ClubID    string `json:"club_id"`
}

// SignupCredentials used when signing up
type SignupCredentials struct {
	Email       string `validate:"required" json:"email"`
	Password    string `validate:"required,gte=4" json:"password"`
	FirstName   string `validate:"required" json:"first_name"`
	LastName    string `validate:"required" json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	// this field is NOT required. If it's present, it will override the
	// `password` field and auto generate a password for the user. Any auto
	// generated password will send an email out to the user so they know
	// what their password is to login. This endpoint is ONLY used on the admin
	// side when creating a `New Golfer`.
	PasswordGen    bool   `json:"password_gen"`
	GolfCourseName string `json:"golf_course_name"`
}

// ResetPasswordRequest request to up password
type ResetPasswordRequest struct {
	ID       string `validate:"required" json:"id"`
	Password string `validate:"required" json:"password"`
}

// ForgotPasswordRequest request to send reset email
type ForgotPasswordRequest struct {
	Email     string `validate:"required" json:"email"`
	CompanyID string `json:"company_id"`
}
