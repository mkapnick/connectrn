package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
)

// queries are written to use sqlx.NamedExec() method. this method maps "db" struct tags with
// the : prefixed names in the values parameter
const (
	// `account` queries
	CreateAccountQuery = `INSERT INTO account (id, email, password, enabled, created_at, updated_at)
VALUES (:id, :email, :password, :enabled, :created_at, :updated_at);`
	UpdateAccountQuery           = `UPDATE account SET email = :email, password = :password WHERE id = :id;`
	FetchAccountQuery            = `SELECT * FROM account WHERE id = $1;`
	FetchAccountByEmailQuery     = `SELECT * FROM account WHERE email = $1;`
	FetchAccountByConditionQuery = `SELECT * FROM account WHERE`
	FetchCompanyQuery            = `SELECT id, subdomain FROM company WHERE id = $1`
	// !!! IMPORTANT !!! In order for me to know the drop rate of users, purging
	// will happen later in a cron. For now, mark the account with enabled = false
	// and set the `deleted_at` date on the account
	DeleteAccountQuery = `UPDATE account SET enabled = :enabled, deleted_at = :deleted_at WHERE id = :id;`
	// DeleteAccountQuery       = `DELETE FROM "user" WHERE id = $1;`

	// `account_role` queries
	CreateAccountRoleQuery      = `INSERT INTO account_role (id, account_id, role_id, authority, company_id, golf_course_id, club_id, created_at, updated_at) VALUES (:id, :account_id, :role_id, :authority, :company_id, :golf_course_id, :club_id, :created_at, :updated_at);`
	FetchAccountRoleByCondition = `SELECT * FROM account_role WHERE`

	// `role` queries
	FetchRoleQuery = `SELECT * FROM role WHERE authority = $1`

	// `password_reset_token` queries
	FetchPasswordResetTokenQuery  = `SELECT * FROM password_reset_token WHERE id = $1`
	CreatePasswordResetTokenQuery = `INSERT INTO password_reset_token (id, company_id, account_id, email, is_used, expires_at, created_at, updated_at) VALUES (:id, :company_id, :account_id, :email, :is_used, :expires_at, :created_at, :updated_at)`
	UpdatePasswordResetTokenQuery = `UPDATE password_reset_token SET is_used = :is_used, updated_at = :updated_at WHERE id = :id`
)

// accountManager is a private implementation of the account.AccountStoreinterface
type accountManager struct {
	// a sqlx database object
	db *sqlx.DB
}

// NewAccountStore returns a postgres db implementation of the account.AccountStore interface
func NewAccountStore(db *sqlx.DB) account.AccountManager {
	return &accountManager{
		db: db,
	}
}

func (s *accountManager) DB() *sqlx.DB {
	return s.db
}

// Create creates a account in a postgres db
func (s *accountManager) CreateAccount(tx *sqlx.Tx, a *account.Account) (*account.Account, error) {
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

func (s *accountManager) UpdateAccount(a *account.Account) (*account.Account, error) {
	// perform update on values we allow to change
	row, err := s.db.NamedExec(UpdateAccountQuery, a)
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

func (s *accountManager) DeleteAccount(a *account.Account) (*account.Account, error) {
	row, err := s.db.NamedExec(DeleteAccountQuery, a)
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

func (s *accountManager) FetchAccount(ID string) (*account.Account, error) {
	var a account.Account

	err := s.db.Get(&a, FetchAccountQuery, ID)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *accountManager) FetchAccountByEmail(email string) (*account.Account, error) {
	var a account.Account

	err := s.db.Get(&a, FetchAccountByEmailQuery, email)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *accountManager) FetchAccountByCondition(whereCondition string) (*account.Account, error) {
	var a account.Account

	query := fmt.Sprintf("%s %s", FetchAccountByConditionQuery, whereCondition)

	err := s.db.Get(&a, query)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *accountManager) FetchCompany(ID string) (*account.Company, error) {
	var c account.Company

	err := s.db.Get(&c, FetchCompanyQuery, ID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *accountManager) CreateAccountRole(tx *sqlx.Tx, ar *account.AccountRole) (*account.AccountRole, error) {
	row, err := tx.NamedExec(CreateAccountRoleQuery, ar)
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

	return ar, nil
}

func (s *accountManager) FetchAccountRoleByCondition(whereCondition string) (*account.AccountRole, error) {
	var ar account.AccountRole

	query := fmt.Sprintf("%s %s", FetchAccountRoleByCondition, whereCondition)

	err := s.db.Get(&ar, query)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (s *accountManager) FetchAllAccountRolesByCondition(whereCondition string) ([]*account.AccountRole, error) {
	var ars []*account.AccountRole

	query := fmt.Sprintf("%s %s", FetchAccountRoleByCondition, whereCondition)

	err := s.db.Select(&ars, query)
	if err != nil {
		return nil, err
	}

	return ars, nil
}

func (s *accountManager) FetchRole(authority string) (*account.Role, error) {
	var r account.Role

	err := s.db.Get(&r, FetchRoleQuery, authority)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Create creates a account in a postgres db
func (s *accountManager) CreatePasswordResetToken(p *account.PasswordResetToken) (*account.PasswordResetToken, error) {
	row, err := s.db.NamedExec(CreatePasswordResetTokenQuery, p)
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

func (s *accountManager) UpdatePasswordResetToken(p *account.PasswordResetToken) (*account.PasswordResetToken, error) {
	// perform update on values we allow to change
	row, err := s.db.NamedExec(UpdatePasswordResetTokenQuery, p)
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

func (s *accountManager) FetchPasswordResetToken(ID string) (*account.PasswordResetToken, error) {
	var p account.PasswordResetToken

	err := s.db.Get(&p, FetchPasswordResetTokenQuery, ID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
