// +build integration

package postgres_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account"
	"gitlab.com/michaelk99/birrdi/api-soa/services/account/postgres"
	"github.com/jmoiron/sqlx"
	"os"
)

const (
	DeleteAccountQuery = `DELETE FROM accounts;`
	CreateAccountQuery = `INSERT INTO accounts (id, email, password, created_at, updated_at, enabled)
							VALUES (:id, :email, :password, :created_at, :updated_at, :enabled);`
	FetchAccountQuery = `SELECT * FROM accounts WHERE id = $1;`
)

func Setup(t *testing.T) (*sqlx.DB, account.AccountStore, func(), func(), *sqlx.Tx) {
	connString := "host=localhost user=postgres dbname=birrdi password='' sslmode=disable"
	if os.Getenv("POSTGRES_CONN_STRING") != "" {
		connString = os.Getenv("POSTGRES_CONN_STRING")
	}
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		t.Fatalf("setup: could not open connection to db: %s", err)
	}

	ds := postgres.NewAccountStore(db)
	tx, _ := db.Beginx()

	return db, ds, func() {
			defer db.Close()
			_, err := db.Exec(DeleteAccountQuery)
			if err != nil {
				t.Fatalf("failed to delete accounts: %v. manual cleanup is necessary", err)
			}
		},
		func() {
			_, err := db.Exec(DeleteAccountQuery)
			if err != nil {
				t.Fatalf("failed to delete accounts: %v. manual cleanup is necessary", err)
			}
		}, tx
}

var AccountTestingTable = []*account.Account{
	&account.Account{
		ID:        uuid.New().String(),
		Email:     "user1@user1.com",
		Password:  "pass",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Enabled:   true,
	},
	&account.Account{
		ID:        uuid.New().String(),
		Email:     "user2@user2.com",
		Password:  "pass",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Enabled:   true,
	},
}

// Test the creation of an account struct into postgres DB
func TestCreateAccount(t *testing.T) {
	_, ds, teardown, _, tx := Setup(t)
	defer teardown()

	for _, a := range AccountTestingTable {
		_, err := ds.Create(tx, a)
		if err != nil {
			t.Fatalf("failed to insert account into DB: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			t.Fatalf("failed to insert account into DB: %v", err)
		}
	}
}

// Test the update of a user struct into postgres DB
func TestUpdateUser(t *testing.T) {
	db, ds, teardown, clearDB, _ := Setup(t)
	defer teardown()

	for _, a := range AccountTestingTable {
		_, err := db.NamedExec(CreateAccountQuery, a)
		if err != nil {
			t.Fatalf("failed to insert account into DB: %v", err)
		}

		a.Email = "change@email.com"
		_, err = ds.Update(a)
		if err != nil {
			t.Fatalf("failed to update account: %v", err)
		}

		// perform fetch and assert
		var aa account.Account
		err = db.Get(&aa, FetchAccountQuery, a.ID)
		if err != nil {
			t.Fatalf("failed to fetch updated user: %v", err)
		}

		assert.Equal(t, "change@email.com", aa.Email)

		clearDB()
	}
}

// Test the delete of a user struct into postgres DB
func TestDeleteUser(t *testing.T) {
	db, ds, teardown, clearDB, _ := Setup(t)
	defer teardown()

	for _, a := range AccountTestingTable {
		_, err := db.NamedExec(CreateAccountQuery, a)
		if err != nil {
			t.Fatalf("failed to insert account into DB: %v", err)
		}

		err = ds.Delete(a.ID)
		if err != nil {
			t.Fatalf("failed to delete account: %v", err)
		}
		clearDB()
	}

	rows, err := db.Queryx("SELECT * FROM accounts")
	if err != nil {
		t.Fatalf("failed to query db: %s", err)
	}

	assert.Equal(t, false, rows.Next())
}

// Test the fetch of a user struct into postgres DB
func TestFetchUserByEmail(t *testing.T) {
	db, ds, teardown, clearDB, _ := Setup(t)
	defer teardown()

	for _, a := range AccountTestingTable {
		_, err := db.NamedExec(CreateAccountQuery, a)
		if err != nil {
			t.Fatalf("failed to insert acc into DB: %v", err)
		}

		aa, err := ds.FetchByEmail(a.Email)

		if err != nil {
			t.Fatalf("failed to get account from DB: %v", err)
		}

		assert.Equal(t, a.Email, aa.Email)
		clearDB()
	}
}

// Test the fetch of a user struct by email from postgres DB
func TestFetchUser(t *testing.T) {
	db, ds, teardown, clearDB, _ := Setup(t)
	defer teardown()

	for _, a := range AccountTestingTable {
		_, err := db.NamedExec(CreateAccountQuery, a)
		if err != nil {
			t.Fatalf("failed to insert acc into DB: %v", err)
		}

		aa, err := ds.Fetch(a.ID)

		if err != nil {
			t.Fatalf("failed to get account from DB: %v", err)
		}

		assert.Equal(t, a.ID, aa.ID)
		clearDB()
	}
}
