package account

import (
	"regexp"
	"time"

	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/connectrn/internal/crypto"
	"gitlab.com/michaelk99/connectrn/internal/transaction"
	"gitlab.com/michaelk99/connectrn/services/profile"
)

// Service is a public interface for implementing our Account service
type Service interface {
	// SignUp creates an account into the backing account store
	SignUp(req SignupCredentials) (*Account, error)
	// LogIn takes account credentials and returns a token if a successful login occurs
	LogIn(ctx context.Context, req AccountCredentials) (token string, err error)
	// Fetch retrieves a user from the backing account store
	Fetch(q IDQuery) (*Account, error)
}

// Service is a private implementation of our account Service
type service struct {
	as AccountStore
	tc TokenCreator
	pc profile.Client
}

// NewService is a constructor for our Account service implementation
func NewService(as AccountStore, tc TokenCreator, pc profile.Client) Service {
	return &service{
		as: as,
		tc: tc,
		pc: pc,
	}
}

// SignUp registers a new user and persists them to the backing account store
func (s *service) SignUp(req SignupCredentials) (*Account, error) {
	// check if account email exists in db. if err is nil an account was found.
	// Right now we ONLY handle `golfer` signups through this endpoint. We do
	// not handle `company` or `club` signups through here. That still happens
	// on the `v1` endpoint from inviting a user.
	// check if account exists in `db`

	// lowercase the string
	req.Email = strings.ToLower(req.Email)

	// email regex
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	isValid := re.MatchString(req.Email)

	// not a valid email address
	if !isValid {
		return nil, ErrInternal{msg: "Invalid email"}
	}

	whereCondition := fmt.Sprintf("restaurant_id IS NULL AND email = '%s'", req.Email)
	_, err := s.as.FetchAccountByCondition(whereCondition)
	if err == nil {
		return nil, ErrUserExists{}
	}

	// create account
	accountID := uuid.New().String()
	ts := time.Now().Format(time.RFC3339)
	a := &Account{
		ID:        accountID,
		Email:     req.Email,
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	// this is only needed for the notification email
	var jwtToken string
	// we need to attach the `profile_id` to the response. The `profile_id`
	// is used from the response when admins create new users
	var dbProfile *profile.Profile

	pass, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, ErrPasswordHash{}
	}

	// set the hashed password to the account struct
	a.Password = pass

	// run the following queries in a DB transaction
	fn := func(tx *sqlx.Tx) error {
		// create account
		_, err := s.as.CreateAccount(tx, a)
		if err != nil {
			return ErrInternal{msg: err.Error()}
		}

		// fashion token [with mock profileID and mock authority]
		token, err := s.tc.Create(a, &profile.Profile{ID: "-1"})

		// set this now, to be used in noitification email if eligible
		jwtToken = token

		if err != nil {
			return ErrCreateToken{err}
		}

		// create profile with mocked token
		dbProfile, err = s.pc.CreateProfile(token, &profile.Profile{
			AccountID: accountID,
		})

		if err != nil {
			return ErrProf{fmt.Sprintf("Could not create initial profile: %s", err)}
		}

		return nil
	}

	err = transaction.Transact(s.as.DB(), fn)
	if err != nil {
		return nil, err
	}

	// we fetch here in order to return the entire Account object
	aa, err := s.as.FetchAccount(accountID)
	if err != nil {
		return nil, ErrInternal{msg: err.Error()}
	}

	aa.ProfileID = dbProfile.ID
	return aa, nil
}

// LogIn authenticates an account's credentials and returns a token if successful.
func (s *service) LogIn(ctx context.Context, req AccountCredentials) (string, error) {
	// lowercase the string
	req.Email = strings.ToLower(req.Email)

	// check if account exists in `db`
	whereCondition := fmt.Sprintf("restaurant_id IS NULL AND email = '%s'", req.Email)
	if req.RestaurantID != "" {
		whereCondition = fmt.Sprintf("restaurant_id = '%s' AND email = '%s'", req.RestaurantID, req.Email)
	}

	// simply look
	a, err := s.as.FetchAccountByCondition(whereCondition)
	if err != nil {
		return "", ErrUserNotFound{}
	}

	// confirm password
	valid := crypto.ValidatePassword(req.Password, a.Password)
	if !valid {
		return "", ErrInvalidLogin{}
	}

	// fashion token [with mock profileID and mock authority]
	// a bit hacky - but the mocked token is used to fetch the user's profile
	token, err := s.tc.Create(a, &profile.Profile{ID: "-1"})
	if err != nil {
		return "", ErrCreateToken{err}
	}

	// get profile
	// !!! IMPORTANT !!!
	// the profile should be created on behalf of the user when they create
	// an account, so this `fetch` should never be an issue
	prof, err := s.pc.GetProfileByAccountID(token, a.ID)
	if err != nil {
		return "", ErrCreateToken{fmt.Errorf("Profile not found for account %s", a.ID)}
	}

	// fashion token again with correct prof ID and authority
	token, err = s.tc.Create(a, prof)

	if err != nil {
		return "", ErrCreateToken{err}
	}

	return token, nil
}

// Fetch retrieves a user from the backing account store. We return an error if
// any issues occurs
func (s *service) Fetch(q IDQuery) (*Account, error) {
	var a *Account
	var err error

	switch q.Type {
	case EmailID:
		a, err = s.as.FetchAccountByEmail(q.Value)
		if err != nil {
			return nil, ErrUserNotFound{}
		}
	case ID:
		a, err = s.as.FetchAccount(q.Value)
		if err != nil {
			return nil, ErrUserNotFound{}
		}
	default:
		return nil, ErrInvalidIDType{}
	}

	return a, nil
}
