package handlers

import (
    "fmt"  // Add this import
    "github.com/gin-gonic/gin"
    "net/http"
    "mobilerecharge/models"
)

func (h *Handler) Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        fmt.Printf("Login binding error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    fmt.Printf("Login attempt with username: %s\n", loginData.Username) // Debug log

    var user models.User
    result := h.DB.Where("username = ? AND password = ?", loginData.Username, loginData.Password).First(&user)
    
    if result.Error != nil {
        fmt.Printf("Login failed for user %s: %v\n", loginData.Username, result.Error)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    fmt.Printf("Login successful for user: %s\n", loginData.Username) // Debug log

    // Set cookie with proper settings
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("logged_in", "true", 3600, "/", "", false, true)

    // Try direct redirect instead of JSON response
    c.Redirect(http.StatusFound, "/")
}