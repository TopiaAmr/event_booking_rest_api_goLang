// Package routes defines the HTTP route handlers and URL patterns for the event booking API.
// It registers all API endpoints with the Gin router and maps them to their handler functions.
package routes

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all API routes with the provided Gin engine.
// It sets up the following endpoints:
//   - GET /events/:id - Get a specific event by ID
//   - GET /events - Get all events
//   - POST /event - Create a new event
func RegisterRoutes(server *gin.Engine) {
	server.GET("/events/:id", getEvent)
	server.GET("/events", getEvents)
	server.POST("/event", createEvent)
}
