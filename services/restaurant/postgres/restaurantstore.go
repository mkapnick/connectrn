package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/connectrn/services/restaurant"
)

// queries are written to use sqlx.NamedExec() method. this method maps "db" struct tags with
// the : prefixed names in the values parameter
const (
	CreateRestaurantQuery          = `INSERT INTO restaurants (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)`
	CreateTableQuery               = `INSERT INTO tables (id, restaurant_id, name, num_seats_available, num_seats_reserved, start_date, created_at, updated_at) VALUES (:id, :restaurant_id, :name, :num_seats_available, :num_seats_reserved, :start_date, :created_at, :updated_at) `
	FetchRestaurantQuery           = `SELECT * FROM restaurants WHERE id = $1`
	FetchAllTablesByConditionQuery = `SELECT * FROM tables WHERE`
)

// restaurantStore is a private implementation of the restaurant.RestaurantStore interface
type restaurantStore struct {
	// a sqlx database object
	db *sqlx.DB
}

// NewRestaurantStore returns a postgres db implementation of the profile.ProfileStore interface
func NewRestaurantStore(db *sqlx.DB) restaurant.RestaurantStore {
	return &restaurantStore{
		db: db,
	}
}

func (s *restaurantStore) CreateRestaurant(r restaurant.Restaurant) (*restaurant.Restaurant, error) {
	// save `user_tee_time_canceled` in the db
	row, err := s.db.NamedExec(CreateRestaurantQuery, r)
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

	return &r, nil
}

func (s *restaurantStore) CreateTable(t restaurant.Table) (*restaurant.Table, error) {
	// save `user_tee_time_canceled` in the db
	row, err := s.db.NamedExec(CreateRestaurantQuery, t)
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

	return &t, nil
}

func (s *restaurantStore) FetchRestaurant(id string) (*restaurant.Restaurant, error) {
	var r restaurant.Restaurant

	err := s.db.Get(&r, FetchRestaurantQuery, id)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *restaurantStore) FetchAllTablesByCondition(whereCondition string) ([]*restaurant.Table, error) {
	var t []*restaurant.Table

	query := fmt.Sprintf("%s %s", FetchAllTablesByConditionQuery, whereCondition)
	err := s.db.Select(&t, query)
	if err != nil {
		return nil, err
	}

	return t, nil
}
