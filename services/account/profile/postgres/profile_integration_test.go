// +build integration

package postgres_test

import (
	"database/sql"
	"gopkg.in/guregu/null.v3"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile"
	"gitlab.com/michaelk99/birrdi/api-soa/services/profile/postgres"
)

const (
	DeleteProfileQuery = `DELETE FROM profiles;`
	CreateProfileQuery = `INSERT INTO profiles (id, account_id, created_at, updated_at) VALUES (:id, :account_id, :created_at, :updated_at);`
	FetchProfileQuery  = `SELECT * FROM profiles WHERE id = $1;`
)

func Setup(t *testing.T) (*sqlx.DB, profile.ProfileStore, func(), func()) {
	connString := "host=localhost user=postgres dbname=birrdi password='' sslmode=disable"
	if os.Getenv("POSTGRES_CONN_STRING") != "" {
		connString = os.Getenv("POSTGRES_CONN_STRING")
	}
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		t.Fatalf("setup: could not open connection to db: %s", err)
	}

	ds := postgres.NewProfileStore(db)

	return db, ds, func() {
			defer db.Close()
			_, err := db.Exec(DeleteProfileQuery)
			if err != nil {
				t.Fatalf("failed to delete profiles: %v. manual cleanup is necessary", err)
			}
		},
		func() {
			_, err := db.Exec(DeleteProfileQuery)
			if err != nil {
				t.Fatalf("failed to delete profiles: %v. manual cleanup is necessary", err)
			}
		}
}

var ProfileTestingTable = []*profile.Profile{
	&profile.Profile{
		ID:        uuid.New().String(),
		AccountID: uuid.New().String(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	},
	&profile.Profile{
		ID:        uuid.New().String(),
		AccountID: uuid.New().String(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	},
}

// Test the creation of an profile struct into postgres DB
func TestCreateProfile(t *testing.T) {
	_, ds, teardown, _ := Setup(t)
	defer teardown()

	for _, p := range ProfileTestingTable {
		_, err := ds.Create(p)
		if err != nil {
			t.Fatalf("failed to insert profile into DB: %v", err)
		}
	}
}

// Test the update of a profile struct into postgres DB
func TestUpdateProfile(t *testing.T) {
	db, ds, teardown, clearDB := Setup(t)
	defer teardown()

	for _, p := range ProfileTestingTable {
		_, err := db.NamedExec(CreateProfileQuery, p)
		if err != nil {
			t.Fatalf("failed to insert profile into DB: %v", err)
		}

		p.Name = sql.NullString{null.String{String: "name-change", Valid: true}}
		_, err = ds.Update(p)
		if err != nil {
			t.Fatalf("failed to update profile: %v", err)
		}

		// perform fetch and assert
		var pp profile.Profile
		err = db.Get(&pp, FetchProfileQuery, p.ID)
		if err != nil {
			t.Fatalf("failed to fetch updated profile: %v", err)
		}

		assert.Equal(t, "name-change", pp.Name.Value)

		clearDB()
	}
}

// Test the delete of a profile struct into postgres DB
func TestDeleteProfile(t *testing.T) {
	db, ds, teardown, clearDB := Setup(t)
	defer teardown()

	for _, p := range ProfileTestingTable {
		_, err := db.NamedExec(CreateProfileQuery, p)
		if err != nil {
			t.Fatalf("failed to insert profile into DB: %v", err)
		}

		err = ds.Delete(p.ID)
		if err != nil {
			t.Fatalf("failed to delete profile: %v", err)
		}
		clearDB()
	}

	rows, err := db.Queryx("SELECT * FROM profile")
	if err != nil {
		t.Fatalf("failed to query db: %s", err)
	}

	assert.Equal(t, false, rows.Next())
}

func TestDeleteProfileByAccountID(t *testing.T) {
	db, ds, teardown, clearDB := Setup(t)
	defer teardown()

	for _, p := range ProfileTestingTable {
		_, err := db.NamedExec(CreateProfileQuery, p)
		if err != nil {
			t.Fatalf("failed to insert profile into DB: %v", err)
		}

		err = ds.DeleteByAccountID(p.AccountID)
		if err != nil {
			t.Fatalf("failed to delete profile: %v", err)
		}
		clearDB()
	}

	rows, err := db.Queryx("SELECT * FROM profile")
	if err != nil {
		t.Fatalf("failed to query db: %s", err)
	}

	assert.Equal(t, false, rows.Next())
}

// Test the fetch of a profile struct by email from postgres DB
func TestFetchProfile(t *testing.T) {
	db, ds, teardown, clearDB := Setup(t)
	defer teardown()

	for _, p := range ProfileTestingTable {
		_, err := db.NamedExec(CreateProfileQuery, p)
		if err != nil {
			t.Fatalf("failed to insert profile into DB: %v", err)
		}

		pp, err := ds.Fetch(p.ID)

		if err != nil {
			t.Fatalf("failed to get profile from DB: %v", err)
		}

		assert.Equal(t, p.ID, pp.ID)
		clearDB()
	}
}

// Test the fetch of a profile struct by email from postgres DB
func TestFetchProfileByAccountID(t *testing.T) {
	db, ds, teardown, clearDB := Setup(t)
	defer teardown()

	for _, p := range ProfileTestingTable {
		_, err := db.NamedExec(CreateProfileQuery, p)
		if err != nil {
			t.Fatalf("failed to insert profile into DB: %v", err)
		}

		pp, err := ds.FetchByAccountID(p.AccountID)

		if err != nil {
			t.Fatalf("failed to get profile from DB: %v", err)
		}

		assert.Equal(t, p.ID, pp.ID)
		clearDB()
	}
}
