// Package main is the entry point for the event booking REST API.
// It initializes the database, sets up the HTTP server, and registers all routes.
package main

import (
	"event_booking_restapi_golang/db"
	"event_booking_restapi_golang/routes"

	"github.com/gin-gonic/gin"
)

// main is the application entry point.
// It initializes the database connection, creates a Gin HTTP server,
// registers all API routes, and starts the server on port 8080.
func main() {
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
