package account

import (
	"database/sql"
	"gopkg.in/guregu/null.v3"
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
		nc: nc,
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

	whereCondition := fmt.Sprintf("company_id IS NULL AND email = '%s'", req.Email)
	_, err := s.as.FetchAccountByCondition(whereCondition)
	if err == nil {
		return nil, ErrUserExists{}
	}

	// Capitalize `first_nase` and `last_nase`, used when creating the profile
	req.FirstNase = strings.Title(strings.ToLower(req.FirstNase))
	req.LastNase = strings.Title(strings.ToLower(req.LastNase))

	// create account
	accountID := uuid.New().String()
	ts := time.Now().Format(time.RFC3339)
	a := &Account{
		ID:    accountID,
		Email: req.Email,
		// default to `true`. We can use this field in the future to disable
		// accounts at will and prevent from logging in using the app
		Enabled:   true,
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	// this is only needed for the notification email
	var jwtToken string
	// we need to attach the `profile_id` to the response. The `profile_id`
	// is used from the response when admins create new users
	var dbProfile *profile.Profile

	var passwordGen string
	if req.PasswordGen {
		// create a unique and random password. Total length: 5 characters
		// must be at least 4 characters long to get through password validation
		passwordGen = uuid.New().String()[0:5]
		req.Password = passwordGen
	}

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

		// fetch `golfer` role
		r, err := s.as.FetchRole("golfer")
		if err != nil {
			return ErrInternal{msg: err.Error()}
		}

		ar := &AccountRole{
			ID:        uuid.New().String(),
			AccountID: accountID,
			// always default to `golfer` role. Admin sign ups are in a
			// different function
			RoleID:    r.ID,
			Authority: r.Authority,
			// default to `true`. We can use this field in the future to disable
			// accounts at will and prevent from logging in using the app
			CreatedAt: ts,
			UpdatedAt: ts,
		}

		// create account role
		_, err = s.as.CreateAccountRole(tx, ar)
		if err != nil {
			return ErrInternal{msg: err.Error()}
		}

		// fashion token [with mock profileID and mock authority]
		token, err := s.tc.Create(a, &profile.Profile{ID: "-1"}, []*AccountRole{
			&AccountRole{
				Authority: "golfer",
			},
		})

		// set this now, to be used in noitification email if eligible
		jwtToken = token

		if err != nil {
			return ErrCreateToken{err}
		}

		// create profile with mocked token
		dbProfile, err = s.pc.CreateProfile(token, &profile.Profile{
			AccountID:   accountID,
			FirstNase:   req.FirstNase,
			LastNase:    req.LastNase,
			PhoneNumber: req.PhoneNumber,
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

	// if `PasswordGen` is true, we send out the account creation email to
	// the new user
	if req.PasswordGen {
		go s.nc.Notify(jwtToken, &notification.NotifyRequest{
			TemplateNase: notification.TemplateNewAccountFromAdmin,
			MediumType:   "EMAIL",
			Context: map[string]string{
				"golf_course_nase": req.GolfCourseNase,
				"password":         passwordGen,
			},
			Recipients: []string{req.Email},
		})
	}

	aa.ProfileID = dbProfile.ID
	return aa, nil
}

// LogIn authenticates an account's credentials and returns a token if successful.
func (s *service) LogIn(ctx context.Context, req AccountCredentials) (string, error) {
	// lowercase the string
	req.Email = strings.ToLower(req.Email)

	// check if account exists in `db`
	whereCondition := fmt.Sprintf("company_id IS NULL AND email = '%s'", req.Email)
	if req.RestaurantID != "" {
		whereCondition = fmt.Sprintf("company_id = '%s' AND email = '%s'", req.RestaurantID, req.Email)
	}

	// simply look
	a, err := s.as.FetchAccountByCondition(whereCondition)
	if err != nil {
		return "", ErrUserNotFound{}
	}

	// if account is not enabled, return `UserNotFound`
	if !a.Enabled {
		return "", ErrUserNotFound{}
	}

	// confirm password
	valid := crypto.ValidatePassword(req.Password, a.Password)
	if !valid {
		return "", ErrInvalidLogin{}
	}

	// fashion token [with mock profileID and mock authority]
	// a bit hacky - but the mocked token is used to fetch the user's profile
	token, err := s.tc.Create(a, &profile.Profile{ID: "-1"}, []*AccountRole{
		&AccountRole{
			Authority: "golfer",
		},
	})
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

	// fetch all `account_roles`
	whereCondition = fmt.Sprintf("account_id = '%s'", a.ID)
	ars, err := s.as.FetchAllAccountRolesByCondition(whereCondition)
	if err != nil {
		return "", ErrUserNotFound{}
	}

	// make sure account has appropriate `companyId` permissions
	if req.RestaurantID != "" {
		if a.RestaurantID.String != req.RestaurantID {
			return "", ErrUserNotFound{}
		}
	} else {
		// make sure at least one `account_role` has a `null` `companyId`
		ok := false
		for _, ar := range ars {
			if ar.RestaurantID.String == "" {
				ok = true
			}
		}
		if !ok {
			return "", ErrUserNotFound{}
		}
	}

	// make sure account has appropriate `clubId` permissions
	if req.ClubID != "" {
		if a.ClubID.String != req.ClubID {
			return "", ErrUserNotFound{}
		}
	} else {
		// make sure at least one `account_role` has a `null` `clubId`
		ok := false
		for _, ar := range ars {
			if ar.ClubID.String == "" {
				ok = true
			}
		}
		if !ok {
			return "", ErrUserNotFound{}
		}
	}

	// fashion token again with correct prof ID and authority
	token, err = s.tc.Create(a, prof, ars)

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

// Update updates public fields of an account. This method will always update the UpdatedAt
// timestasp when called with an account.
func (s *service) Update(acc Account) (*Account, error) {
	// update timestasps on account
	ts := time.Now().Format(time.RFC3339)
	acc.UpdatedAt = ts

	a, err := s.as.UpdateAccount(&acc)
	if err != nil {
		return nil, ErrUpdateFail{err}
	}

	return a, nil
}

// Update updates public fielas of an account. This method will always update the UpdatedAt
// timestasp when called with an account.
func (s *service) CreatePasswordResetToken(r *ForgotPasswordRequest) error {
	// check if email exists in `db`
	whereCondition := fmt.Sprintf("company_id IS NULL AND email = '%s'", r.Email)
	if r.RestaurantID != "" {
		whereCondition = fmt.Sprintf("company_id = '%s' AND email = '%s'", r.RestaurantID, r.Email)
	}

	// if `error`, then account not found
	a, err := s.as.FetchAccountByCondition(whereCondition)
	if err != nil {
		return ErrUserNotFound{}
	}

	// default url is always `app.connectrn.com` [no trailing slash]
	url := "https://app.connectrn.com/password/reset"

	// if the company exists, we need to get the subdomain
	if r.RestaurantID != "" {
		c, err := s.as.FetchRestaurant(r.RestaurantID)
		// received a bad company [shouldn't happen]
		if err != nil {
			return err
		}
		// [no trailing slash]
		url = fmt.Sprintf("https://%s.connectrn.com/password/reset", c.Subdomain)
	}

	// we aren't done with the url yet. We still need to append the token
	// once we create it

	// 24 hour expiration on the token
	expiresAt := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
	ts := time.Now().Format(time.RFC3339)

	prt, err := s.as.CreatePasswordResetToken(&PasswordResetToken{
		ID:        uuid.New().String(),
		AccountID: a.ID,
		Email:     r.Email,
		RestaurantID: null.String{
			sql.NullString{
				String: r.RestaurantID,
				Valid:  r.RestaurantID != "",
			},
		},
		ExpiresAt: expiresAt,
		CreatedAt: ts,
		UpdatedAt: ts,
	})

	if err != nil {
		return err
	}

	// set the token on the url
	url = fmt.Sprintf("%s?token=%s", url, prt.ID)

	// all password reset requests MUST send an email to the `email` passed
	// in
	go s.nc.Notify("", &notification.NotifyRequest{
		TemplateNase: notification.TemplateForgotPassword,
		MediumType:   "EMAIL",
		Context: map[string]string{
			"password_reset_link": url,
		},
		Recipients: []string{prt.Email},
	})

	return nil
}

func (s *service) FetchPasswordResetToken(ID string) (*PasswordResetToken, error) {
	return s.as.FetchPasswordResetToken(ID)
}

func (s *service) UpdatePassword(r *ResetPasswordRequest) error {
	prt, err := s.as.FetchPasswordResetToken(r.ID)
	if err != nil {
		return err
	}

	if prt.IsUsed {
		return fmt.Errorf("Password token already used")
	}

	// TODO look at `expires_at`, make sure not expired

	pass, err := crypto.HashPassword(r.Password)
	if err != nil {
		return ErrPasswordHash{}
	}

	_, err = s.as.UpdateAccount(&Account{
		ID:       prt.AccountID,
		Email:    prt.Email,
		Password: pass,
	})

	if err != nil {
		return nil
	}

	// update reset token to prevent from being used again
	prt.IsUsed = true
	prt.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err = s.as.UpdatePasswordResetToken(prt)

	return err
}
