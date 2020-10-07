package account

import (
	"github.com/jmoiron/sqlx"
)

// AccountManager is responsible for managing accounts
type AccountManager interface {
	AccountStore
	AccountRoleStore
	PasswordResetTokenStore
}

// AccountStore interface
type AccountStore interface {
	CreateAccount(*sqlx.Tx, *Account) (*Account, error)
	UpdateAccount(*Account) (*Account, error)
	FetchAccount(ID string) (*Account, error)
	FetchAccountByEmail(email string) (*Account, error)
	FetchAccountByCondition(whereCondition string) (*Account, error)
	DeleteAccount(*Account) (*Account, error)

	// For transaction purposes only
	DB() *sqlx.DB
}

// AccountRoleStore interface
type AccountRoleStore interface {
	CreateAccountRole(*sqlx.Tx, *AccountRole) (*AccountRole, error)
	FetchAccountRoleByCondition(whereCondition string) (*AccountRole, error)
	FetchAllAccountRolesByCondition(whereCondition string) ([]*AccountRole, error)
	FetchRole(authority string) (*Role, error)
}

// PasswordResetTokenStore interface
type PasswordResetTokenStore interface {
	CreatePasswordResetToken(*PasswordResetToken) (*PasswordResetToken, error)
	UpdatePasswordResetToken(*PasswordResetToken) (*PasswordResetToken, error)
	FetchPasswordResetToken(ID string) (*PasswordResetToken, error)
	// need to get the subdomain of the `company` to construct the appropriate
	// link
	FetchCompany(ID string) (*Company, error)
}
