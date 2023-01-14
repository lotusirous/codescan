package db

import (
	"context"
	"database/sql"
	"fmt"
)

// DB wraps the database and provides a helper.
type DB struct {
	conn *sql.DB
}

// Exec execute the query in the context.
func (db *DB) Exec(ctx context.Context, query string, args ...any) error {
	r, err := db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	return nil
}

// ExecTX wraps the execution in the transaction.
func (db *DB) ExecTX(ctx context.Context, query string, args ...any) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	r, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	return tx.Commit()
}

// Tx makes a new transaction
func (db *DB) Tx() (*sql.Tx, error) {
	return db.conn.Begin()
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.conn.QueryRowContext(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

func (db *DB) Close() error {
	return db.conn.Close()
}
