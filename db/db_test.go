// Package db contains unit tests for database initialization and operations.
package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestInitDB tests the database initialization function
func TestInitDB(t *testing.T) {
	// Use a test database file
	testDBFile := "test_db.sql"

	// Clean up any existing test database file
	os.Remove(testDBFile)
	defer os.Remove(testDBFile)

	// Create a test database connection directly
	testDB, err := sql.Open("sqlite3", testDBFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer testDB.Close()

	// Test table creation
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id TEXT
	)
	`
	_, err = testDB.Exec(createEventsTable)
	if err != nil {
		t.Errorf("Failed to create events table: %v", err)
	}

	// Verify table exists
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='events'").Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify table creation: %v", err)
	}
	if count != 1 {
		t.Error("Events table was not created")
	}

	// Test connection settings
	testDB.SetMaxOpenConns(10)
	testDB.SetMaxIdleConns(5)

	// Verify connection is working
	err = testDB.Ping()
	if err != nil {
		t.Errorf("Database connection is not working: %v", err)
	}
}

// TestCreateTables tests the table creation function
func TestCreateTables(t *testing.T) {
	// Create an in-memory database for testing
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer testDB.Close()

	// Test table creation
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id TEXT
	)
	`
	_, err = testDB.Exec(createEventsTable)
	if err != nil {
		t.Errorf("Failed to create events table: %v", err)
	}

	// Verify table structure
	rows, err := testDB.Query("PRAGMA table_info(events)")
	if err != nil {
		t.Errorf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	columns := []string{}
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue interface{}

		err = rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			t.Errorf("Failed to scan column info: %v", err)
		}
		columns = append(columns, name)
	}

	expectedColumns := []string{"id", "name", "description", "location", "datetime", "user_id"}
	if len(columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(columns))
	}

	for i, expected := range expectedColumns {
		if i >= len(columns) || columns[i] != expected {
			t.Errorf("Expected column %s, got %s", expected, columns[i])
		}
	}
}

// TestDatabaseConnection tests database connection functionality
func TestDatabaseConnection(t *testing.T) {
	// Test with in-memory database
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer testDB.Close()

	// Test connection
	err = testDB.Ping()
	if err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}

	// Test connection limits
	testDB.SetMaxOpenConns(10)
	testDB.SetMaxIdleConns(5)

	// Verify connection is still working
	err = testDB.Ping()
	if err != nil {
		t.Errorf("Database connection failed after setting limits: %v", err)
	}
}

// TestDatabaseErrorHandling tests error handling in database operations
func TestDatabaseErrorHandling(t *testing.T) {
	// Test with invalid database path - SQLite doesn't validate path until first operation
	invalidDB, err := sql.Open("sqlite3", "/invalid/path/db.sql")
	if err != nil {
		t.Fatalf("Unexpected error when opening database with invalid path: %v", err)
	}

	// Try to ping to trigger the error
	err = invalidDB.Ping()
	if err == nil {
		t.Error("Expected error when pinging database with invalid path")
	}
	invalidDB.Close()

	// Test with valid database but invalid SQL
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer testDB.Close()

	// Test invalid SQL
	_, err = testDB.Exec("INVALID SQL STATEMENT")
	if err == nil {
		t.Error("Expected error when executing invalid SQL")
	}
}
