package handlers

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "mobilerecharge/models"
)

func (h *Handler) Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    err := h.db.Collection("users").FindOne(
        context.Background(),
        bson.M{
            "username": loginData.Username,
            "password": loginData.Password,
        },
    ).Decode(&user)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    c.SetCookie("logged_in", "true", 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}