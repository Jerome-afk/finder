package server

import "github.com/gin-gonic/gin"

// Initialize the server
func InitServer() *gin.Engine {
	router := gin.Default()

	// Initialise logs
	initLogger()

	return router
}
