package server

import (
	"log"

	"github/finderr/backend/db"

	"github.com/gin-gonic/gin"
)

// Initialize the server
func InitServer() *gin.Engine {
	router := gin.Default()

	// Initialise logs
	initLogger()

	// Initialise the database
	database, err := db.InitDB("./finder.db")
	if err != nil {
		LogEvent("Server", "Failed to initialise database")
		log.Fatalf("Failed to initialize database: %v", err)
	}
	LogEvent("Server", "Database initialised successfully")
	defer database.Close()

	// Run migrations
	err = db.RunMigrations(database)
	if err != nil {
		LogEvent("Server", "Failed to run migration")
		log.Fatalf("Failed to run migrations: %v", err)
	}
	LogEvent("Server", "Migration run successfully")

	return router
}
