package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/connectrn/services/profile"
)

// queries are written to use sqlx.NamedExec() method. this method maps "db" struct tags with
// the : prefixed names in the values parameter
const (
	CreateProfileQuery           = `INSERT INTO profiles (id, account_id, name, created_at, updated_at) VALUES (:id, :account_id, :name, :created_at, :updated_at);`
	FetchProfileQuery            = `SELECT * FROM profiles WHERE id = $1;`
	FetchProfileByAccountIDQuery = `SELECT * FROM profiles WHERE account_id = $1`
)

// profileStore is a private implementation of the profile.ProfileStore interface
type profileStore struct {
	// a sqlx database object
	db *sqlx.DB
}

// NewProfileStore returns a postgres db implementation of the profile.ProfileStore interface
func NewProfileStore(db *sqlx.DB) profile.ProfileStore {
	return &profileStore{
		db: db,
	}
}

// Create creates a profile in a postgres db
func (s *profileStore) Create(p *profile.Profile) (*profile.Profile, error) {
	row, err := s.db.NamedExec(CreateProfileQuery, p)
	if err != nil {
		return nil, err
	}

	i, err := row.RowsAffected()
	switch {
	case i <= 0:
		return nil, fmt.Errorf("%d rows affected by update", i)
	case err != nil:
		return nil, err
	}

	return p, nil
}

func (s *profileStore) Fetch(ID string) (*profile.Profile, error) {
	var p profile.Profile

	err := s.db.Get(&p, FetchProfileQuery, ID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *profileStore) FetchByAccountID(accountID string) (*profile.Profile, error) {
	var p profile.Profile

	err := s.db.Get(&p, FetchProfileByAccountIDQuery, accountID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
