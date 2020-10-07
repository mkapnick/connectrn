package account

import (
	"github.com/jmoiron/sqlx"
)

// AccountStore interface
type AccountStore interface {
	CreateAccount(*sqlx.Tx, *Account) (*Account, error)
	FetchAccount(ID string) (*Account, error)
	FetchAccountByEmail(email string) (*Account, error)
	FetchAccountByCondition(whereCondition string) (*Account, error)

	// For transaction purposes only
	DB() *sqlx.DB
}
