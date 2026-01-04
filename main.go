package main

import (
	"event_booking_restapi_golang/db"
	"event_booking_restapi_golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
