package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//Initialise the database
func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil,err
	}

	// Check if database is connected
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}