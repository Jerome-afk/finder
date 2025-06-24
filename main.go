package main

import (
        "log"

        "finderr/internal/api"
        "finderr/internal/config"
        "finderr/internal/database"
        "finderr/internal/middleware"
        "finderr/internal/services"

        "github.com/gin-gonic/gin"
        "github.com/joho/godotenv"
)

func main() {
        // Load environment variables
        if err := godotenv.Load(); err != nil {
                log.Println("No .env file found, using system environment variables")
        }

        // Initialize configuration
        cfg := config.Load()
        log.Printf("Database path: %s", cfg.DatabasePath)

        // Initialize database
        db, err := database.Initialize(cfg.DatabasePath)
        if err != nil {
                log.Fatal("Failed to initialize database:", err)
        }
        defer db.Close()
        log.Println("Database connection established")

        // Run migrations
        if err := database.RunMigrations(cfg.DatabasePath); err != nil {
                log.Printf("Migration warning: %v", err)
                // Don't fail on migration errors in case tables already exist
        }
        log.Println("Database migrations completed")

        // Initialize services
        userService := services.NewUserService(db)
        movieService := services.NewMovieService(cfg.TMDBAPIKey, cfg.OMDBAPIKey)
        watchlistService := services.NewWatchlistService(db)
        log.Println("Services initialized")

        // Initialize Gin router
        router := gin.Default()

        // Add middleware
        router.Use(middleware.CORS())
        router.Use(middleware.SessionMiddleware(cfg.SessionSecret))

        // Serve static files
        router.Static("/static", "./web/static")
        router.LoadHTMLGlob("web/templates/*")

        // Initialize API routes
        api.SetupRoutes(router, userService, movieService, watchlistService)
        log.Println("Routes configured")

        // Start server
        port := cfg.Port
        log.Printf("Starting server on 0.0.0.0:%s", port)
        log.Printf("Server ready at http://localhost:%s", port)
        
        if err := router.Run("0.0.0.0:" + port); err != nil {
                log.Fatal("Failed to start server:", err)
        }
}