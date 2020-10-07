package transaction

import (
	"github.com/jmoiron/sqlx"
)

// Transact db transaction for all or nothing queries
func Transact(db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
