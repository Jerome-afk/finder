package api

import (
        "net/http"

        "github.com/gin-gonic/gin"
)

type WebHandler struct{}

func NewWebHandler() *WebHandler {
        return &WebHandler{}
}

func (h *WebHandler) HomePage(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
                "title": "FindeRR - Media Tracker",
        })
}

func (h *WebHandler) LoginPage(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", gin.H{
                "title": "Login - FindeRR",
        })
}

func (h *WebHandler) RegisterPage(c *gin.Context) {
        c.HTML(http.StatusOK, "register.html", gin.H{
                "title": "Register - FindeRR",
        })
}

func (h *WebHandler) SearchPage(c *gin.Context) {
        query := c.Query("q")
        c.HTML(http.StatusOK, "search.html", gin.H{
                "title": "Search - FindeRR",
                "query": query,
        })
}

func (h *WebHandler) DiscoverPage(c *gin.Context) {
        c.HTML(http.StatusOK, "discover.html", gin.H{
                "title": "Discover - FindeRR",
        })
}

func (h *WebHandler) MovieDetailPage(c *gin.Context) {
        movieID := c.Param("id")
        c.HTML(http.StatusOK, "movie_detail.html", gin.H{
                "title":   "Movie Details - FindeRR",
                "movieID": movieID,
        })
}

func (h *WebHandler) TVDetailPage(c *gin.Context) {
        tvID := c.Param("id")
        c.HTML(http.StatusOK, "tv_detail.html", gin.H{
                "title": "TV Show Details - FindeRR",
                "tvID":  tvID,
        })
}

func (h *WebHandler) DashboardPage(c *gin.Context) {
        c.HTML(http.StatusOK, "dashboard.html", gin.H{
                "title": "Dashboard - FindeRR",
        })
}

func (h *WebHandler) WatchlistPage(c *gin.Context) {
        c.HTML(http.StatusOK, "watchlist.html", gin.H{
                "title": "My Watchlist - FindeRR",
        })
}

func (h *WebHandler) StatsPage(c *gin.Context) {
        c.HTML(http.StatusOK, "stats.html", gin.H{
                "title": "My Stats - FindeRR",
        })
}