package api

import (
	"net/http"
	"strconv"

	"finderr/internal/services"

	"github.com/gin-gonic/gin"
)

type WatchlistHandler struct {
	watchlistService *services.WatchlistService
	userService      *services.UserService
}

func NewWatchlistHandler(watchlistService *services.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{watchlistService: watchlistService}
}

type AddToWatchlistRequest struct {
	MediaID    int    `json:"media_id" binding:"required"`
	MediaType  string `json:"media_type" binding:"required,oneof=movie tv"`
	Title      string `json:"title" binding:"required"`
	PosterPath string `json:"poster_path"`
}

type UpdateWatchStatusRequest struct {
	IsWatched bool `json:"is_watched"`
}

type RateMediaRequest struct {
	Rating int `json:"rating" binding:"required,min=1,max=10"`
}

func (h *WatchlistHandler) GetWatchlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	mediaType := c.Query("type")
	var isWatched *bool
	if watchedStr := c.Query("watched"); watchedStr != "" {
		if watched, err := strconv.ParseBool(watchedStr); err == nil {
			isWatched = &watched
		}
	}

	items, err := h.watchlistService.GetUserWatchlist(userID, mediaType, isWatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *WatchlistHandler) AddToWatchlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	var req AddToWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.watchlistService.AddToWatchlist(userID, req.MediaID, req.MediaType, req.Title, req.PosterPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to watchlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Added to watchlist successfully",
		"item":    item,
	})
}

func (h *WatchlistHandler) RemoveFromWatchlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	mediaIDStr := c.Param("mediaId")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	mediaType := c.Param("mediaType")
	if mediaType != "movie" && mediaType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media type"})
		return
	}

	err = h.watchlistService.RemoveFromWatchlist(userID, mediaID, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from watchlist successfully"})
}

func (h *WatchlistHandler) UpdateWatchStatus(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	mediaIDStr := c.Param("mediaId")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	mediaType := c.Param("mediaType")
	if mediaType != "movie" && mediaType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media type"})
		return
	}

	var req UpdateWatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.watchlistService.UpdateWatchStatus(userID, mediaID, mediaType, req.IsWatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update watch status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch status updated successfully"})
}

func (h *WatchlistHandler) RateMedia(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	mediaIDStr := c.Param("mediaId")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	mediaType := c.Param("mediaType")
	if mediaType != "movie" && mediaType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media type"})
		return
	}

	var req RateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.watchlistService.RateMedia(userID, mediaID, mediaType, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media rated successfully"})
}

func (h *WatchlistHandler) CheckWatchlist(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	mediaIDStr := c.Param("mediaId")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	mediaType := c.Param("mediaType")
	if mediaType != "movie" && mediaType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media type"})
		return
	}

	inWatchlist, err := h.watchlistService.IsInWatchlist(userID, mediaID, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"in_watchlist": inWatchlist})
}

func (h *WatchlistHandler) GetRecommendations(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	recommendations, err := h.watchlistService.GetRecommendations(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}

func (h *WatchlistHandler) GetUserStats(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userID := userIDStr.(string)

	stats, err := h.watchlistService.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}