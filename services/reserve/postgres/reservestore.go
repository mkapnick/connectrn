package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/birrdi/api-soa/services/reserve"
)

// queries are written to use sqlx.NamedExec() method. this method maps "db" struct tags with
// the : prefixed names in the values parameter
const (
	FetchAllReservesByConditionQuery = `SELECT * FROM golf_course_reserve_date WHERE`
)

// reserveStore is a private implementation of the reserve.BlackoutDateStore interface
type reserveStore struct {
	// a sqlx database object
	db *sqlx.DB
}

// NewReserveStore returns a postgres db implementation of the profile.ProfileStore interface
func NewReserveStore(db *sqlx.DB) reserve.ReserveStore {
	return &reserveStore{
		db: db,
	}
}

func (s *reserveStore) FetchAllByCondition(whereCondition string) ([]*reserve.Reserve, error) {
	var b []*reserve.Reserve

	query := fmt.Sprintf("%s %s", FetchAllReservesByConditionQuery, whereCondition)
	err := s.db.Select(&b, query)
	if err != nil {
		fmt.Printf("BAD %s", err)
		return nil, err
	}

	return b, nil
}
