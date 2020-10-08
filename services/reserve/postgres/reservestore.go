package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/connectrn/services/reserve"
)

// queries are written to use sqlx.NamedExec() method. this method maps "db" struct tags with
// the : prefixed names in the values parameter
const (
	CreateUserReservationQuery         = `INSERT INTO user_reservations (id, restaurant_id, table_id, profile_id, num_seats, start_date, created_at, updated_at) VALUES (:id, :restaurant_id, :table_id, :profile_id, :num_seats, :start_date, :created_at, :updated_at)`
	CreateUserReservationCanceledQuery = `INSERT INTO user_reservations_canceled (id, restaurant_id, table_id, profile_id, num_seats, start_date, created_at, updated_at) VALUES (:id, :restaurant_id, :table_id, :profile_id, :num_seats, :start_date, :created_at, :updated_at)`
	DeleteUserReservationQuery         = `DELETE FROM user_reservations WHERE id = $1`

	// table queries
	FetchTableQuery  = `SELECT * FROM tables WHERE restaurant_id = $1 AND id = $2`
	UpdateTableQuery = `UPDATE tables SET num_seats_reserved = :num_seats_reserved WHERE id = :id`
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

// used to create a proper db transaction
func (s *reserveStore) DB() *sqlx.DB {
	return s.db
}

func (s *reserveStore) FetchTable(restaurantID string, tableID string) (*reserve.Table, error) {
	var t reserve.Table

	err := s.db.Get(&t, FetchTableQuery, restaurantID, tableID)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *reserveStore) CreateUserReservation(tx *sqlx.Tx, r reserve.UserReservation) (*reserve.UserReservation, error) {
	// save `user_tee_time_canceled` in the db
	row, err := tx.NamedExec(CreateUserReservationQuery, r)
	if err != nil {
		return nil, err
	}

	i, err := row.RowsAffected()
	switch {
	case i <= 0:
		return nil, fmt.Errorf("%d rows affected by create", i)
	case err != nil:
		return nil, err
	}

	return &r, nil
}

func (s *reserveStore) CreateUserReservationCanceled(tx *sqlx.Tx, r reserve.UserReservationCanceled) (*reserve.UserReservationCanceled, error) {
	// save `user_tee_time_canceled` in the db
	row, err := tx.NamedExec(CreateUserReservationCanceledQuery, r)
	if err != nil {
		return nil, err
	}

	i, err := row.RowsAffected()
	switch {
	case i <= 0:
		return nil, fmt.Errorf("%d rows affected by create", i)
	case err != nil:
		return nil, err
	}

	return &r, nil
}

func (s *reserveStore) UpdateTable(tx *sqlx.Tx, t reserve.Table) (*reserve.Table, error) {
	// save `user_tee_time_canceled` in the db
	row, err := tx.NamedExec(UpdateTableQuery, t)
	if err != nil {
		return nil, err
	}

	i, err := row.RowsAffected()
	switch {
	case i <= 0:
		return nil, fmt.Errorf("%d rows affected by create", i)
	case err != nil:
		return nil, err
	}

	return &t, nil
}

func (s *reserveStore) DeleteUserReservation(tx *sqlx.Tx, ID string) error {
	res, err := tx.Exec(DeleteUserReservationQuery, ID)
	if err != nil {
		return err
	}

	i, err := res.RowsAffected()
	switch {
	case i <= 0:
		return fmt.Errorf("%d rows affected by delete", i)
	case err != nil:
		return err
	}

	return nil
}
