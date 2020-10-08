package token

// Validator is an interface for a generic token store
type Validator interface {
	Validate(tokenString string) (*Session, error)
}

// Session is the session derived from the token
type Session struct {
	AccountID    string `json:"account_id"`
	ProfileID    string `json:"profile_id"`
	Email        string `json:"email"`
	RestaurantID string `json:"restaurant_id"`
}
