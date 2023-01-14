package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Connect to a database and verify with a ping.
func Connect(datasource string, maxOpenConnections int) (*sql.DB, error) {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, err
	}
	if err := pingDatabase(db); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConnections)

	return db, nil

}

// helper function to ping the database with backoff to ensure
// a connection can be established before we proceed with the
// database setup.
func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}
