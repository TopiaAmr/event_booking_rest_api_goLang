// Package testutils provides common utilities for testing the event booking API.
package testutils

import (
	"database/sql"
	"event_booking_restapi_golang/db"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TestDB holds the test database connection
type TestDB struct {
	DB         *sql.DB
	OriginalDB *sql.DB
}

// SetupTestDatabase creates a fresh in-memory SQLite database for testing
// and returns a TestDB struct that can be used to clean up after tests
func SetupTestDatabase(t *testing.T) *TestDB {
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create events table for testing
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id TEXT
	)
	`
	_, err = testDB.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Store original DB and replace with test DB
	originalDB := db.DB
	db.DB = testDB

	return &TestDB{
		DB:         testDB,
		OriginalDB: originalDB,
	}
}

// Cleanup restores the original database connection and closes the test database
func (tdb *TestDB) Cleanup() {
	if tdb.OriginalDB != nil {
		db.DB = tdb.OriginalDB
	}
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// CreateTestEvent creates a test event with default values
func CreateTestEvent() map[string]interface{} {
	return map[string]interface{}{
		"title":       "Test Event",
		"description": "Test Description",
		"location":    "Test Location",
		"datetime":    time.Now().Add(time.Hour).Format(time.RFC3339),
	}
}

// CreateTestEventWithCustomData creates a test event with custom values
func CreateTestEventWithCustomData(title, description, location string, datetime time.Time) map[string]interface{} {
	return map[string]interface{}{
		"title":       title,
		"description": description,
		"location":    location,
		"datetime":    datetime.Format(time.RFC3339),
	}
}

// AssertDatabaseCount verifies the number of records in a table
func AssertDatabaseCount(t *testing.T, db *sql.DB, table string, expectedCount int) {
	var count int
	query := "SELECT COUNT(*) FROM " + table
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		t.Errorf("Failed to count records in %s: %v", table, err)
	}
	if count != expectedCount {
		t.Errorf("Expected %d records in %s, got %d", expectedCount, table, count)
	}
}

// AssertEventExists verifies that an event exists in the database
func AssertEventExists(t *testing.T, db *sql.DB, title string) {
	var count int
	query := "SELECT COUNT(*) FROM events WHERE name = ?"
	err := db.QueryRow(query, title).Scan(&count)
	if err != nil {
		t.Errorf("Failed to check if event exists: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected event '%s' to exist, but found %d occurrences", title, count)
	}
}

// AssertEventNotExists verifies that an event does not exist in the database
func AssertEventNotExists(t *testing.T, db *sql.DB, title string) {
	var count int
	query := "SELECT COUNT(*) FROM events WHERE name = ?"
	err := db.QueryRow(query, title).Scan(&count)
	if err != nil {
		t.Errorf("Failed to check if event exists: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected event '%s' to not exist, but found %d occurrences", title, count)
	}
}

// InsertTestEvent inserts a test event directly into the database
func InsertTestEvent(t *testing.T, db *sql.DB, title, description, location string, datetime time.Time, userID string) string {
	query := `
	INSERT INTO events (id, name, description, location, datetime, user_id)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	eventID := "test-event-" + time.Now().Format("20060102150405")
	_, err := db.Exec(query, eventID, title, description, location, datetime, userID)
	if err != nil {
		t.Fatalf("Failed to insert test event: %v", err)
	}

	return eventID
}

// GetEventByID retrieves an event from the database by ID
func GetEventByID(t *testing.T, db *sql.DB, id string) map[string]interface{} {
	query := "SELECT id, name, description, location, datetime, user_id FROM events WHERE id = ?"
	row := db.QueryRow(query, id)

	var eventID, name, description, location, userID string
	var datetime time.Time

	err := row.Scan(&eventID, &name, &description, &location, &datetime, &userID)
	if err != nil {
		t.Errorf("Failed to get event by ID: %v", err)
		return nil
	}

	return map[string]interface{}{
		"id":          eventID,
		"title":       name,
		"description": description,
		"location":    location,
		"datetime":    datetime,
		"user_id":     userID,
	}
}

// ClearEventsTable removes all events from the events table
func ClearEventsTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM events")
	if err != nil {
		t.Errorf("Failed to clear events table: %v", err)
	}
}
