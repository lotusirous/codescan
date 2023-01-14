package dbtest

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lotusirous/codescan/store/db"
)

// Connect opens a new test database connection.
func Connect() (*db.DB, error) {
	var datasource string
	datasource = os.Getenv("DATABASE_DATASOURCE")
	return db.Connect(datasource, 0)
}

var noContext = context.TODO()

// Reset resets the database state.
func Reset(d *db.DB) error {
	tx, err := d.Tx()
	if err != nil {
		return err
	}
	tx.Exec("DELETE FROM repositories")
	tx.Exec("DELETE FROM scans")
	tx.Exec("DELETE FROM scan_results")
	return tx.Commit()
}

// Disconnect closes the database connection.
func Disconnect(d *db.DB) error {
	return d.Close()
}
