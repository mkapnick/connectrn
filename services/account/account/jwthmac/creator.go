package jwthmac

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
)

// CustomClaim ability to add custom claim data here
type CustomClaim struct {
	jwt.StandardClaims
	AccountID    string                 `json:"account_id"`
	ProfileID    string                 `json:"profile_id"`
	CompanyID    string                 `json:"company_id"`
	ClubID       string                 `json:"club_id"`
	AccountRoles []*account.AccountRole `json:"account_roles"`

	// combine both `account` and `profile` data within the jwt
	Email       string      `json:"email"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Preferences interface{} `json:"preferences"`
}

// creator creates hmac signed jwt tokens
type creator struct {
	// secret used to sign hmac jwt tokens
	secret string
	// the expiration offset from current time. this value will be added to time.Now()
	expiration time.Duration
	// the issuer of the hmac signed jwt
	issuer string
}

// NewCreator is a constructor for our hmac jwt token creator
func NewCreator(secret, issuer string, exp time.Duration) account.TokenCreator {
	return &creator{
		secret:     secret,
		expiration: exp,
		issuer:     issuer,
	}
}

// Create method implements the token.Creator interface. Takes a subject string
// and uses the struct's configured state to fashion an hmac signed token
func (c *creator) Create(acc *account.Account, prof *profile.Profile, ars []*account.AccountRole) (string, error) {
	ts := time.Now()
	iat := ts.Unix()
	exp := ts.Add(c.expiration).Unix()
	jti := uuid.New().String()

	// create a new jwt
	claims := CustomClaim{
		jwt.StandardClaims{
			ExpiresAt: exp,
			Id:        jti,
			IssuedAt:  iat,
			Issuer:    c.issuer,
			Subject:   acc.ID,
		},
		acc.ID,
		prof.ID,
		acc.CompanyID.String,
		acc.ClubID.String,
		ars,
		acc.Email,
		prof.FirstName,
		prof.LastName,
		prof.Preferences,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(c.secret))

	if err != nil {
		return "", err
	}

	return s, nil
}
