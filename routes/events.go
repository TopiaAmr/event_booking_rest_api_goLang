// Package routes contains the HTTP handler functions for the event booking API endpoints.
// It implements the business logic for handling HTTP requests and responses.
package routes

import (
	"event_booking_restapi_golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getEvents handles GET requests to /events endpoint.
// It retrieves all events from the database and returns them as JSON.
// Returns HTTP 500 if there's an error fetching events, otherwise HTTP 200 with events data.
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err, "where": "couldn't fetch events"})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

// getEvent handles GET requests to /events/:id endpoint.
// It retrieves a specific event by its ID from the database.
// Returns HTTP 404 if the event is not found, otherwise HTTP 302 with the event data.
func getEvent(c *gin.Context) {
	id, _ := c.Params.Get("id")
	event, err := models.GetEventById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"event": event,
	})

}

// createEvent handles POST requests to /event endpoint.
// It creates a new event from the JSON request body and saves it to the database.
// Returns HTTP 400 if the request is invalid or save fails, otherwise HTTP 201 with the created event.
func createEvent(context *gin.Context) {
	var newEvent models.Event
	err := context.ShouldBindJSON(&newEvent)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{"message": "something went wrong", "error": err},
		)
		return
	}
	newEvent.ID = uuid.NewString()
	newEvent.UserID = uuid.NewString()
	err = newEvent.Save()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{"message": "something went wrong", "error": err},
		)
		return
	}
	context.JSON(
		http.StatusCreated,
		gin.H{"message": "A new event has been created successfully", "event": newEvent},
	)
}
