package token

import (
	"gitlab.com/michaelk99/connectrn/services/account"
)

// Validator is an interface for a generic token store
type Validator interface {
	Validate(tokenString string) (*Session, error)
}

// Session is the session derived from the token
type Session struct {
	AccountID    string                 `json:"account_id"`
	ProfileID    string                 `json:"profile_id"`
	Email        string                 `json:"email"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	CompanyID    string                 `json:"company_id"`
	ClubID       string                 `json:"club_id"`
	AccountRoles []*account.AccountRole `json:"account_roles"`
}
