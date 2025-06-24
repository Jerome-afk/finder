package api

import (
	"finderr/internal/middleware"
	"finderr/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	userService *services.UserService,
	movieService *services.MovieService,
	watchlistService *services.WatchlistService,
) {
	// Initialize handlers
	authHandler := NewAuthHandler(userService)
	movieHandler := NewMovieHandler(movieService)
	watchlistHandler := NewWatchlistHandler(watchlistService)
	webHandler := NewWebHandler()

	// Public routes
	router.GET("/", webHandler.HomePage)
	router.GET("/login", webHandler.LoginPage)
	router.GET("/register", webHandler.RegisterPage)
	router.GET("/search", webHandler.SearchPage)
	router.GET("/movie/:id", webHandler.MovieDetailPage)
	router.GET("/tv/:id", webHandler.TVDetailPage)

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/user", middleware.AuthRequired(), authHandler.GetCurrentUser)
	}

	// API routes
	api := router.Group("/api")
	api.Use(middleware.RateLimiter())
	{
		// Public movie/TV endpoints
		api.GET("/search/movies", movieHandler.SearchMovies)
		api.GET("/search/tv", movieHandler.SearchTVShows)
		api.GET("/search", movieHandler.Search)
		api.GET("/movie/:id", movieHandler.GetMovieDetails)
		api.GET("/tv/:id", movieHandler.GetTVDetails)
		api.GET("/trending", movieHandler.GetTrending)
		api.GET("/movies/genre/:id", movieHandler.GetMoviesByGenre)
		api.GET("/tv/genre/:id", movieHandler.GetTVByGenre)

		// Protected watchlist endpoints
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/watchlist", watchlistHandler.GetWatchlist)
			protected.POST("/watchlist", watchlistHandler.AddToWatchlist)
			protected.DELETE("/watchlist/:mediaId/:mediaType", watchlistHandler.RemoveFromWatchlist)
			protected.PUT("/watchlist/:mediaId/:mediaType/watched", watchlistHandler.UpdateWatchStatus)
			protected.PUT("/watchlist/:mediaId/:mediaType/rating", watchlistHandler.RateMedia)
			protected.GET("/watchlist/check/:mediaId/:mediaType", watchlistHandler.CheckWatchlist)
			protected.GET("/recommendations", watchlistHandler.GetRecommendations)
			protected.GET("/stats", watchlistHandler.GetUserStats)
		}
	}

	// Dashboard (protected)
	dashboard := router.Group("/dashboard")
	dashboard.Use(middleware.AuthRequired())
	{
		dashboard.GET("/", webHandler.DashboardPage)
		dashboard.GET("/watchlist", webHandler.WatchlistPage)
		dashboard.GET("/stats", webHandler.StatsPage)
	}
}
