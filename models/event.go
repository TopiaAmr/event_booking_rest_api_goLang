package models

import (
	"errors"
	"event_booking_restapi_golang/db"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          string
	Title       string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      string
}

var events = []Event{}

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
