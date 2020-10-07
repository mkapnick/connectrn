package profile

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/michaelk99/birrdi/api-soa/services/square"
)

// Service is a public interface for implementing our Profile service
type Service interface {
	// Create creates a profile
	Create(prof Profile) (*Profile, error)
	// Fetch retrieves a profile from the db
	Fetch(query IDQuery) (*Profile, error)
	// Search retrieves all profiles from the db that match
	Search(jwt string, query IDQuery, golfCourseID string) ([]*Profile, error)
	// Update updates a profile in the db
	Update(prof Profile) (*Profile, error)
	// Delete deletes a profile from the db
	Delete(query IDQuery) error
}

// service is a private implementation of our profile service
type service struct {
	ds ProfileStore
	sc square.Client
}

// NewService is a constructor for our Profile service implementation
func NewService(ds ProfileStore, sc square.Client) Service {
	return &service{
		ds: ds,
		sc: sc,
	}
}

func (s *service) Create(prof Profile) (*Profile, error) {
	// fetch by account ID, see if profile already exists
	// 1-1 association between profile and account
	p, _ := s.ds.FetchByAccountID(prof.AccountID)
	if p != nil {
		return nil, ErrProfileExists{}
	}

	ts := time.Now().Format(time.RFC3339)
	prof.CreatedAt = ts
	prof.UpdatedAt = ts
	prof.ID = uuid.New().String()

	_, err := s.ds.Create(&prof)
	if err != nil {
		return nil, ErrProfileCreate{}
	}

	return &prof, nil
}

func (s *service) Search(jwt string, query IDQuery, golfCourseID string) ([]*Profile, error) {
	var profs []*Profile

	switch query.Type {

	case FromLoyalty:
		whereCondition := fmt.Sprintf("(TRIM(account.email) ILIKE '%s%%' OR TRIM(profile.first_name) ILIKE '%s%%' OR TRIM(profile.last_name) ILIKE '%s%%' OR CONCAT(profile.first_name, ' ', profile.last_name) ILIKE '%s%%') AND (account.company_id IS NULL and account.club_id IS NULL) LIMIT 3", query.Value, query.Value, query.Value, query.Value)

		profs, _ = s.ds.FetchAllByCondition(whereCondition)
		return profs, nil

	case AdminCheckout:

		whereCondition := fmt.Sprintf("(TRIM(account.email) ILIKE '%s%%' OR TRIM(profile.first_name) ILIKE '%s%%' OR TRIM(profile.last_name) ILIKE '%s%%' OR CONCAT(profile.first_name, ' ', profile.last_name) ILIKE '%s%%') AND (account.company_id IS NULL and account.club_id IS NULL) LIMIT 3", query.Value, query.Value, query.Value, query.Value)

		profs, _ = s.ds.FetchAllByCondition(whereCondition)

		if len(profs) > 0 {
			return profs, nil
		}

		// if we cannot find account in Birrdi, search Square
		resp, err := s.sc.SearchCustomerFuzzy(jwt, golfCourseID, query.Value)
		if err != nil {
			return profs, nil
		}

		// no customers found
		if len(resp.Customers) == 0 {
			return profs, nil
		}

		// found, get into profile format
		for _, c := range resp.Customers {
			profs = append(profs, &Profile{
				ID:        "",
				Email:     c.Email,
				FirstName: c.GivenName,
				LastName:  c.FamilyName,
			})
		}
	}

	return profs, nil
}

func (s *service) Fetch(query IDQuery) (*Profile, error) {
	var prof *Profile
	var err error

	switch query.Type {
	case ID:
		prof, err = s.ds.Fetch(query.Value)
		if err != nil {
			return nil, ErrProfileNotFound{}
		}
	case AccountID:
		prof, err = s.ds.FetchByAccountID(query.Value)
		if err != nil {
			return nil, ErrProfileNotFound{}
		}
	default:
		return nil, fmt.Errorf("Invalid query type")
	}

	return prof, nil
}

func (s *service) Update(prof Profile) (*Profile, error) {
	ts := time.Now().Format(time.RFC3339)
	prof.UpdatedAt = ts

	_, err := s.ds.Update(&prof)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}

func (s *service) Delete(query IDQuery) error {
	switch query.Type {
	case ID:
		err := s.ds.Delete(query.Value)
		if err != nil {
			return err
		}
	case AccountID:
		err := s.ds.DeleteByAccountID(query.Value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid query type")
	}

	return nil
}
