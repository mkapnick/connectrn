package account

import (
	"gopkg.in/guregu/null.v3"
)

// Account is a retrieved and account
type Account struct {
	ID         string      `json:"id" db:"id"`
	CompanyID  null.String `json:"company_id" db:"company_id"`
	ClubID     null.String `json:"club_id" db:"club_id"`
	Email      string      `json:"email" db:"email"`
	Password   string      `json:"-" db:"password"`
	Enabled    bool        `json:"enabled" db:"enabled"`
	InviteCode null.String `json:"invite_code" db:"invite_code"`
	InvitedAt  null.String `json:"invited_at" db:"invited_at"`
	CreatedAt  string      `json:"created_at" db:"created_at"`
	UpdatedAt  string      `json:"updated_at" db:"updated_at"`

	// this is ONLY used to attach the `profile_id` to the account object.
	// this field does NOT exist as a db field
	ProfileID string `json:"profile_id"`
}

// AccountRole is a retrieved account role
type AccountRole struct {
	ID           string      `json:"id" db:"id"`
	AccountID    string      `json:"account_id" db:"account_id"`
	GolfCourseID null.String `json:"golf_course_id" db:"golf_course_id"`
	CompanyID    null.String `json:"company_id" db:"company_id"`
	ClubID       null.String `json:"club_id" db:"club_id"`
	RoleID       string      `json:"role_id" db:"role_id"`
	Authority    string      `json:"authority" db:"authority"`
	CreatedAt    string      `json:"created_at" db:"created_at"`
	UpdatedAt    string      `json:"updated_at" db:"updated_at"`
}

// Role is a retrieved role
type Role struct {
	ID           string      `json:"id" db:"id"`
	Authority    string      `json:"authority" db:"authority"`
	InviteHeader null.String `json:"invite_header" db:"invite_header"`
	InviteBody   null.String `json:"invite_body" db:"invite_body"`
	CreatedAt    string      `json:"created_at" db:"created_at"`
	UpdatedAt    string      `json:"updated_at" db:"updated_at"`
}

// PasswordResetToken is a retrieved reset token
type PasswordResetToken struct {
	ID        string      `json:"id" db:"id"`
	AccountID string      `json:"account_id" db:"account_id"`
	CompanyID null.String `json:"company_id" db:"company_id"`
	Email     string      `json:"email" db:"email"`
	IsUsed    bool        `json:"is_used" db:"is_used"`
	ExpiresAt string      `json:"expires_at" db:"expires_at"`
	CreatedAt string      `json:"created_at" db:"created_at"`
	UpdatedAt string      `json:"updated_at" db:"updated_at"`
}

// Company retrieves a company subdomain
type Company struct {
	ID        string `json:"id" db:"id"`
	Subdomain string `json:"subdomain" db:"subdomain"`
}
