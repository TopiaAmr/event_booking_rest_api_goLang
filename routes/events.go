package routes

import (
	"event_booking_restapi_golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
