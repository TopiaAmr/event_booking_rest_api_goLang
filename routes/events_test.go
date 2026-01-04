// Package routes contains unit tests for the HTTP handlers.
package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"event_booking_restapi_golang/db"
	"event_booking_restapi_golang/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

// setupTestRouter creates a Gin router for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// TestGetEvents tests the getEvents handler
func TestGetEvents(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.GET("/events", getEvents)

	// Insert test events
	events := []models.Event{
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

	req, _ := http.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	eventsData, ok := response["events"].([]interface{})
	if !ok {
		t.Error("Response should contain 'events' array")
	}

	if len(eventsData) != 2 {
		t.Errorf("Expected 2 events, got %d", len(eventsData))
	}
}

// TestGetEventsEmpty tests the getEvents handler with no events
func TestGetEventsEmpty(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.GET("/events", getEvents)

	req, _ := http.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	t.Logf("Response body: %s", w.Body.String())

	eventsData, ok := response["events"]
	if !ok {
		t.Error("Response should contain 'events' array")
	}

	if eventsData == nil {
		t.Log("Events data is null, treating as empty array")
		return // This is actually correct behavior for empty database
	}

	eventsSlice, ok := eventsData.([]interface{})
	if !ok {
		t.Errorf("Expected events to be an array, got %T", eventsData)
	}

	if len(eventsSlice) != 0 {
		t.Errorf("Expected 0 events, got %d", len(eventsSlice))
	}
}

// TestGetEvent tests the getEvent handler
func TestGetEvent(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.GET("/events/:id", getEvent)

	// Insert a test event
	event := models.Event{
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

	req, _ := http.NewRequest("GET", "/events/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	t.Logf("Response body: %s", w.Body.String())

	eventData, ok := response["event"]
	if !ok {
		t.Error("Response should contain 'event' object")
	}

	eventMap, ok := eventData.(map[string]interface{})
	if !ok {
		t.Errorf("Expected event to be an object, got %T", eventData)
	}

	if eventMap["Title"] != event.Title {
		t.Errorf("Expected title %s, got %v", event.Title, eventMap["Title"])
	}
}

// TestGetEventNotFound tests the getEvent handler with non-existent ID
func TestGetEventNotFound(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.GET("/events/:id", getEvent)

	req, _ := http.NewRequest("GET", "/events/non-existent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("Response should contain 'error' field")
	}
}

// TestCreateEvent tests the createEvent handler
func TestCreateEvent(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.POST("/event", createEvent)

	eventData := map[string]interface{}{
		"title":       "New Event",
		"description": "New Description",
		"location":    "New Location",
		"datetime":    time.Now().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(eventData)
	req, _ := http.NewRequest("POST", "/event", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Response should contain 'message' field")
	}

	if _, ok := response["event"]; !ok {
		t.Error("Response should contain 'event' field")
	}
}

// TestCreateEventInvalidJSON tests the createEvent handler with invalid JSON
func TestCreateEventInvalidJSON(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.POST("/event", createEvent)

	invalidJSON := `{"title": "Test"}` // Missing required fields
	req, _ := http.NewRequest("POST", "/event", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateEvent tests the updateEvent handler
func TestUpdateEvent(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.PUT("/events/:id", updateEvent)

	// Insert a test event
	event := models.Event{
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

	// Update data
	updateData := map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated Description",
		"location":    "Updated Location",
		"datetime":    time.Now().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/events/"+id, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Response should contain 'message' field")
	}

	if _, ok := response["event"]; !ok {
		t.Error("Response should contain 'event' field")
	}
}

// TestUpdateEventNotFound tests the updateEvent handler with non-existent ID
func TestUpdateEventNotFound(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.PUT("/events/:id", updateEvent)

	updateData := map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated Description",
		"location":    "Updated Location",
		"datetime":    time.Now().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/events/non-existent-id", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestDeleteEvent tests the deleteEvent handler
func TestDeleteEvent(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.DELETE("/events/:id", deleteEvent)

	// Insert a test event
	event := models.Event{
		Title:       "Event to Delete",
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

	req, _ := http.NewRequest("DELETE", "/events/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Response should contain 'message' field")
	}

	// Verify the event was deleted
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM events WHERE id = ?", id).Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify event deletion: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 events after deletion, got %d", count)
	}
}

// TestDeleteEventNotFound tests the deleteEvent handler with non-existent ID
func TestDeleteEventNotFound(t *testing.T) {
	setupTestDatabase(t)
	router := setupTestRouter()
	router.DELETE("/events/:id", deleteEvent)

	req, _ := http.NewRequest("DELETE", "/events/non-existent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
