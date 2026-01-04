// Package routes defines the HTTP route handlers and URL patterns for the event booking API.
// It registers all API endpoints with the Gin router and maps them to their handler functions.
package routes

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all API routes with the provided Gin engine.
// It sets up the following endpoints:
//   - GET /events/:id - Get a specific event by ID
//   - GET /events - Get all events
//   - POST /event - Create a new event
//   - PUT /events/:id - Update an existing event
//   - DELETE /events/:id - Delete an event
func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)
	server.POST("/event", createEvent)
	server.PUT("/events/:id", updateEvent)
	server.GET("/events/:id", getEvent)
	server.DELETE("/events/:id", deleteEvent)
}
