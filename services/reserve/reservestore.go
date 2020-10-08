package reserve

import (
	"github.com/jmoiron/sqlx"
)

// ReserveStore interface
type ReserveStore interface {
	// fetches table for data validation
	FetchTable(restaurantID string, tableID string) (*Table, error)
	FetchUserReservation(ID string) (*UserReservation, error)
	CreateUserReservation(tx *sqlx.Tx, r UserReservation) (*UserReservation, error)
	CreateUserReservationCanceled(tx *sqlx.Tx, r UserReservation) (*UserReservation, error)
	// updates `num_spots_reserved`
	UpdateTable(tx *sqlx.Tx, t Table) (*Table, error)
	DeleteUserReservation(tx *sqlx.Tx, r UserReservation) error
	// For transaction purposes only
	DB() *sqlx.DB
}
