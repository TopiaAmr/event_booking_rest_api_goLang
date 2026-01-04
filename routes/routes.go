package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events/:id", getEvent)
	server.GET("/events", getEvents)
	server.POST("/event", createEvent)
}
