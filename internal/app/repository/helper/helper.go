package helper

import (
	"database/sql"
	"fmt"
)

func ExecWithTx[T any](db *sql.DB, fn func(tx *sql.Tx) (T, error)) (T, error) {
	var zero T

	tx, err := db.Begin()
	if err != nil {
		return zero, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	result, err := fn(tx)
	if err != nil {
		return result, err
	}

	return result, nil
}
