// Package models contains unit tests for the Event model and its database operations.
package models

import (
	"database/sql"
	"event_booking_restapi_golang/db"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var testDB *sql.DB

// setupTestDatabase creates a fresh in-memory SQLite database for testing
func setupTestDatabase(t *testing.T) {
	var err error
	testDB, err = sql.Open("sqlite3", ":memory:")
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

	// Replace the global DB with test DB
	originalDB := db.DB
	db.DB = testDB
	t.Cleanup(func() {
		db.DB = originalDB
		testDB.Close()
	})
}

// TestEvent_Save tests the Save method of the Event model
func TestEvent_Save(t *testing.T) {
	setupTestDatabase(t)

	event := Event{
		Title:       "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now(),
		UserID:      "test-user-123",
	}

	err := event.Save()
	if err != nil {
		t.Errorf("Failed to save event: %v", err)
	}

	// Verify the event was saved
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM events WHERE name = ?", event.Title).Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify event was saved: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 event to be saved, got %d", count)
	}
}

// TestGetAllEvents tests the GetAllEvents function
func TestGetAllEvents(t *testing.T) {
	setupTestDatabase(t)

	// Insert test events
	events := []Event{
		{
			Title:       "Event 1",
			Description: "Description 1",
			Location:    "Location 1",
			DateTime:    time.Now(),
			UserID:      "user1",
		},
		{
			Title:       "Event 2",
			Description: "Description 2",
			Location:    "Location 2",
			DateTime:    time.Now().Add(time.Hour),
			UserID:      "user2",
		},
	}

	for _, event := range events {
		err := event.Save()
		if err != nil {
			t.Fatalf("Failed to insert test event: %v", err)
		}
	}

	retrievedEvents, err := GetAllEvents()
	if err != nil {
		t.Errorf("Failed to get all events: %v", err)
	}

	if len(retrievedEvents) != 2 {
		t.Errorf("Expected 2 events, got %d", len(retrievedEvents))
	}
}

// TestGetEventById tests the GetEventById function
func TestGetEventById(t *testing.T) {
	setupTestDatabase(t)

	// Insert a test event
	event := Event{
		Title:       "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now(),
		UserID:      "test-user-123",
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to save test event: %v", err)
	}

	// Get the event ID from the database
	var id string
	err = testDB.QueryRow("SELECT id FROM events WHERE name = ?", event.Title).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to get event ID: %v", err)
	}

	// Test retrieving the event
	retrievedEvent, err := GetEventById(id)
	if err != nil {
		t.Errorf("Failed to get event by ID: %v", err)
	}

	if retrievedEvent.Title != event.Title {
		t.Errorf("Expected title %s, got %s", event.Title, retrievedEvent.Title)
	}

	// Test with non-existent ID
	_, err = GetEventById("non-existent-id")
	if err == nil {
		t.Error("Expected error when getting non-existent event")
	}
}

// TestEvent_Update tests the Update method of the Event model
func TestEvent_Update(t *testing.T) {
	setupTestDatabase(t)

	// Insert a test event
	event := Event{
		Title:       "Original Title",
		Description: "Original Description",
		Location:    "Original Location",
		DateTime:    time.Now(),
		UserID:      "test-user-123",
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to save test event: %v", err)
	}

	// Get the event ID
	var id string
	err = testDB.QueryRow("SELECT id FROM events WHERE name = ?", event.Title).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to get event ID: %v", err)
	}

	// Update the event
	event.ID = id
	event.Title = "Updated Title"
	event.Description = "Updated Description"
	event.Location = "Updated Location"

	err = event.Update()
	if err != nil {
		t.Errorf("Failed to update event: %v", err)
	}

	// Verify the update
	var title, description, location string
	err = testDB.QueryRow("SELECT name, description, location FROM events WHERE id = ?", id).Scan(&title, &description, &location)
	if err != nil {
		t.Errorf("Failed to verify event update: %v", err)
	}

	if title != "Updated Title" || description != "Updated Description" || location != "Updated Location" {
		t.Error("Event was not updated correctly")
	}
}

// TestEvent_Delete tests the Delete method of the Event model
func TestEvent_Delete(t *testing.T) {
	setupTestDatabase(t)

	// Insert a test event
	event := Event{
		Title:       "Test Event to Delete",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now(),
		UserID:      "test-user-123",
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to save test event: %v", err)
	}

	// Get the event ID
	var id string
	err = testDB.QueryRow("SELECT id FROM events WHERE name = ?", event.Title).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to get event ID: %v", err)
	}

	// Delete the event
	event.ID = id
	err = event.Delete()
	if err != nil {
		t.Errorf("Failed to delete event: %v", err)
	}

	// Verify the deletion
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM events WHERE id = ?", id).Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify event deletion: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 events after deletion, got %d", count)
	}
}

// TestGetEventsByUserId tests the GetEventsByUserId function
func TestGetEventsByUserId(t *testing.T) {
	setupTestDatabase(t)

	userId := "test-user-123"

	// Insert test events for the same user
	events := []Event{
		{
			Title:       "User Event 1",
			Description: "Description 1",
			Location:    "Location 1",
			DateTime:    time.Now(),
			UserID:      userId,
		},
		{
			Title:       "User Event 2",
			Description: "Description 2",
			Location:    "Location 2",
			DateTime:    time.Now().Add(time.Hour),
			UserID:      userId,
		},
		{
			Title:       "Other User Event",
			Description: "Other Description",
			Location:    "Other Location",
			DateTime:    time.Now().Add(2 * time.Hour),
			UserID:      "other-user-456",
		},
	}

	for _, event := range events {
		err := event.Save()
		if err != nil {
			t.Fatalf("Failed to insert test event: %v", err)
		}
	}

	// Get events for the specific user
	userEvents, err := GetEventsByUserId(userId)
	if err != nil {
		t.Errorf("Failed to get events by user ID: %v", err)
	}

	if len(userEvents) != 2 {
		t.Errorf("Expected 2 events for user %s, got %d", userId, len(userEvents))
	}

	// Verify all returned events belong to the correct user
	for _, event := range userEvents {
		if event.UserID != userId {
			t.Errorf("Expected user ID %s, got %s", userId, event.UserID)
		}
	}
}

// TestEventValidation tests the validation tags on the Event struct
func TestEventValidation(t *testing.T) {
	// This test would require additional validation logic in the Save method
	// For now, we'll test the basic structure
	event := Event{
		Title:       "",
		Description: "",
		Location:    "",
		DateTime:    time.Time{},
		UserID:      "",
	}

	// The Save method should handle validation through binding tags
	// This is a placeholder for when validation is implemented
	_ = event
}
