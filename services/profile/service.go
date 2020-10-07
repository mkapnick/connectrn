package profile

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service is a public interface for implementing our Profile service
type Service interface {
	// Create creates a profile
	Create(prof Profile) (*Profile, error)
	// Fetch retrieves a profile from the db
	Fetch(query IDQuery) (*Profile, error)
}

// service is a private implementation of our profile service
type service struct {
	ds ProfileStore
}

// NewService is a constructor for our Profile service implementation
func NewService(ds ProfileStore) Service {
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
