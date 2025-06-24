package main

import (
	"github/finderr/backend/server"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Initialise .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	if len(os.Args) > 1 {
		log.Fatal("Too many arguments parsed\n ......[USAGE]: go run . ")
	}

	router := server.InitServer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8787"
	}

	log.Printf("Server starting at: http://localhost:%s....", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
	server.LogEvent("Server", "Server started successfully")
}