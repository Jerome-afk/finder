package config

import (
        "os"
)

type Config struct {
        DatabasePath  string
        TMDBAPIKey    string
        OMDBAPIKey    string
        SessionSecret string
        Port          string
}

func Load() *Config {
        // Set default values for missing secrets to allow basic functionality
        tmdbKey := os.Getenv("TMDB_API_KEY")
        if tmdbKey == "" {
                tmdbKey = "demo" // Will show error message to user about missing key
        }
        
        omdbKey := os.Getenv("OMDB_API_KEY")
        if omdbKey == "" {
                omdbKey = "demo" // Optional, OMDB features will be disabled
        }
        
        sessionSecret := os.Getenv("SESSION_SECRET")
        if sessionSecret == "" {
                sessionSecret = "demo-session-secret-please-set-real-secret-in-production"
        }
        
        return &Config{
                DatabasePath:  getEnvWithDefault("DATABASE_PATH", "./data/finderr.db"),
                TMDBAPIKey:    tmdbKey,
                OMDBAPIKey:    omdbKey,
                SessionSecret: sessionSecret,
                Port:          getEnvWithDefault("PORT", "8080"),
        }
}

func getEnvWithDefault(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        return defaultValue
}