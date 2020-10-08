package account

// AccountCredentials is used when requesting in some nature
// the access to an account.
type AccountCredentials struct {
	Email     string `validate:"required" json:"email"`
	Password  string `validate:"required,gte=4" json:"password"`
	RestaurantID string `json:"restaurant_id"`
}

// SignupCredentials used when signing up
type SignupCredentials struct {
	Email       string `validate:"required" json:"email"`
	Password    string `validate:"required,gte=4" json:"password"`
}
