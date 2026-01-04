// Package models defines the data structures and database operations for events.
// It provides the Event model and functions for CRUD operations on events.
package models

import (
	"errors"
	"event_booking_restapi_golang/db"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event represents an event in the system with all its properties.
// It includes basic event information like title, description, location,
// as well as metadata like ID, date/time, and user ID.
type Event struct {
	ID          string    // Unique identifier for the event
	Title       string    `binding:"required"` // Event title (required)
	Description string    `binding:"required"` // Event description (required)
	Location    string    `binding:"required"` // Event location (required)
	DateTime    time.Time `binding:"required"` // Event date and time (required)
	UserID      string    // ID of the user who created the event
}

// events is a slice used to store events in memory (currently unused in database operations)
var events = []Event{}

// Save persists the Event to the database.
// It generates a new UUID for the event and inserts it into the events table.
// Returns an error if the database operation fails.
func (e Event) Save() error {
	q := `
	INSERT INTO events (id, name,description,datetime,user_id,location)
	VALUES (?,?,?,?,?,?)
	`
	stmt, err := db.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uuid.NewString(), e.Title, e.Description, e.DateTime, e.UserID, e.Location)
	if err != nil {
		return err
	}

	return nil
}

// GetAllEvents retrieves all events from the database.
// Returns a slice of Event objects and any error encountered during the query.
func GetAllEvents() ([]Event, error) {
	q := `SELECT * FROM events`
	rows, err := db.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// GetEventById retrieves a single event from the database by its ID.
// Returns the Event object if found, otherwise returns an empty Event and an error.
func GetEventById(id string) (Event, error) {
	q := "SELECT * FROM events where id=?"
	row := db.DB.QueryRow(q, id)
	var event Event

	err := row.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return Event{}, errors.New(fmt.Sprint("Couldn't find an event with the ID of", id))
	}

	if event.ID == "" {
		return Event{}, errors.New(fmt.Sprint("Couldn't find an event with the ID of", id))
	}

	return event, nil
}
