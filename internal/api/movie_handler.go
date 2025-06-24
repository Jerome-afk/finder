package api

import (
	"net/http"
	"strconv"

	"finderr/internal/services"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieService *services.MovieService
}

func NewMovieHandler(movieService *services.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

func (h *MovieHandler) SearchMovies(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	results, err := h.movieService.SearchMovies(query, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search movies"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *MovieHandler) SearchTVShows(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	results, err := h.movieService.SearchTVShows(query, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search TV shows"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *MovieHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Search both movies and TV shows concurrently
	moviesChan := make(chan interface{})
	tvChan := make(chan interface{})

	go func() {
		if results, err := h.movieService.SearchMovies(query, page); err == nil {
			moviesChan <- results
		} else {
			moviesChan <- nil
		}
	}()

	go func() {
		if results, err := h.movieService.SearchTVShows(query, page); err == nil {
			tvChan <- results
		} else {
			tvChan <- nil
		}
	}()

	movies := <-moviesChan
	tvShows := <-tvChan

	response := gin.H{
		"query": query,
		"page":  page,
	}

	if movies != nil {
		response["movies"] = movies
	}
	if tvShows != nil {
		response["tv_shows"] = tvShows
	}

	c.JSON(http.StatusOK, response)
}

func (h *MovieHandler) GetMovieDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := h.movieService.GetMovieDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movie details"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) GetTVDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TV show ID"})
		return
	}

	tvShow, err := h.movieService.GetTVShowDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get TV show details"})
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

func (h *MovieHandler) GetTrending(c *gin.Context) {
	content, err := h.movieService.GetTrendingContent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trending content"})
		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *MovieHandler) GetMoviesByGenre(c *gin.Context) {
	idStr := c.Param("id")
	genreID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	results, err := h.movieService.GetMoviesByGenre(genreID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies by genre"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *MovieHandler) GetTVByGenre(c *gin.Context) {
	idStr := c.Param("id")
	genreID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	results, err := h.movieService.GetTVShowsByGenre(genreID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get TV shows by genre"})
		return
	}

	c.JSON(http.StatusOK, results)
}
