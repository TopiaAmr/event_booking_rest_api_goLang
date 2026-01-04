// Package db handles database initialization and operations for the event booking API.
// It provides functions to initialize the SQLite database connection and create necessary tables.
package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection pool used throughout the application.
var DB *sql.DB

// InitDB initializes the SQLite database connection and configures connection settings.
// It opens a connection to "db.sql", sets connection limits, and creates required tables.
// Panics if the database connection fails.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "db.sql")

	if err != nil {
		log.Fatal("Couldn't init DB ", err)
		panic(1)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

// createTables creates the necessary database tables for the application.
// Currently creates the events table if it doesn't exist.
// Panics if table creation fails.
func createTables() {
	createEventsTable := `
		CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id int
		)
		`
	_, err := DB.Exec(createEventsTable)
	if err != nil {
		log.Fatal("Couldn't create events table ", err)
		panic(1)
	}
}
