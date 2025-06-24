package api

import (
        "net/http"
        "strings"

        "finderr/internal/services"

        "github.com/gin-contrib/sessions"
        "github.com/gin-gonic/gin"
)

type AuthHandler struct {
        userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
        return &AuthHandler{userService: userService}
}

type RegisterRequest struct {
        Username string `json:"username" binding:"required,min=3,max=50"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
        var req RegisterRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        // Clean inputs
        req.Username = strings.TrimSpace(req.Username)
        req.Email = strings.TrimSpace(strings.ToLower(req.Email))

        // Create user
        user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
        if err != nil {
                if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
                        c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
                        return
                }
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
                return
        }

        // Create session
        session := sessions.Default(c)
        session.Set("user_id", user.ID)
        session.Save()

        c.JSON(http.StatusCreated, gin.H{
                "message": "User created successfully",
                "user": gin.H{
                        "id":       user.ID,
                        "username": user.Username,
                        "email":    user.Email,
                },
        })
}

func (h *AuthHandler) Login(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        // Clean inputs
        req.Email = strings.TrimSpace(strings.ToLower(req.Email))

        // Get user by email
        user, err := h.userService.GetUserByEmail(req.Email)
        if err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
                return
        }

        // Validate password
        if err := h.userService.ValidatePassword(user.Password, req.Password); err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
                return
        }

        // Create session
        session := sessions.Default(c)
        session.Set("user_id", user.ID)
        session.Save()

        c.JSON(http.StatusOK, gin.H{
                "message": "Login successful",
                "user": gin.H{
                        "id":       user.ID,
                        "username": user.Username,
                        "email":    user.Email,
                },
        })
}

func (h *AuthHandler) Logout(c *gin.Context) {
        session := sessions.Default(c)
        session.Clear()
        session.Save()

        c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
        userIDStr, exists := c.Get("user_id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
                return
        }

        userID := userIDStr.(string)

        user, err := h.userService.GetUserByID(userID)
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
                return
        }

        // Note: Stats are now handled by watchlist service
        var stats interface{} = nil

        c.JSON(http.StatusOK, gin.H{
                "user": gin.H{
                        "id":       user.ID,
                        "username": user.Username,
                        "email":    user.Email,
                },
                "stats": stats,
        })
}