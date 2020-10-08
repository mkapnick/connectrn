package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/connectrn/services/account"
)

// queries are written to use sqlx.NamedExec() method. this method maps "db" struct tags with
// the : prefixed names in the values parameter
const (
	// `account` queries
	CreateAccountQuery = `INSERT INTO accounts (id, email, password, restaurant_id, created_at, updated_at) VALUES (:id, :email, :password, :restaurant_id, :created_at, :updated_at);`
	FetchAccountQuery            = `SELECT * FROM accounts WHERE id = $1;`
	FetchAccountByEmailQuery     = `SELECT * FROM accounts WHERE email = $1;`
	FetchAccountByConditionQuery = `SELECT * FROM accounts WHERE`
	FetchRestaurantQuery         = `SELECT id, name FROM restaurant WHERE id = $1`
)

// accountStore is a private implementation of the account.AccountStoreinterface
type accountStore struct {
	// a sqlx database object
	db *sqlx.DB
}

// NewAccountStore returns a postgres db implementation of the account.AccountStore interface
func NewAccountStore(db *sqlx.DB) account.AccountStore {
	return &accountStore{
		db: db,
	}
}

func (s *accountStore) DB() *sqlx.DB {
	return s.db
}

// Create creates a account in a postgres db
func (s *accountStore) CreateAccount(tx *sqlx.Tx, a *account.Account) (*account.Account, error) {
	row, err := tx.NamedExec(CreateAccountQuery, a)
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

	return a, nil
}

func (s *accountStore) FetchAccount(ID string) (*account.Account, error) {
	var a account.Account

	err := s.db.Get(&a, FetchAccountQuery, ID)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *accountStore) FetchAccountByEmail(email string) (*account.Account, error) {
	var a account.Account

	err := s.db.Get(&a, FetchAccountByEmailQuery, email)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *accountStore) FetchAccountByCondition(whereCondition string) (*account.Account, error) {
	var a account.Account

	query := fmt.Sprintf("%s %s", FetchAccountByConditionQuery, whereCondition)

	err := s.db.Get(&a, query)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *accountStore) FetchRestaurant(ID string) (*account.Restaurant, error) {
	var c account.Restaurant

	err := s.db.Get(&c, FetchRestaurantQuery, ID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
