package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

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
