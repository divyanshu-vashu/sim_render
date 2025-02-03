package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "mobilerecharge/models"
    "mobilerecharge/handlers"
    "mobilerecharge/config"
    "mobilerecharge/services"
    "time"
)

// Add this middleware function
func authRequired(c *gin.Context) {
    loggedIn, _ := c.Cookie("logged_in")
    if loggedIn != "true" {
        c.Redirect(http.StatusFound, "/login")
        c.Abort()
        return
    }
    c.Next()
}

func main() {
    // Initialize MongoDB connection
    config.ConnectDB()
    
    // Disconnect when the main function exits
    defer func() {
        if err := config.MongoClient.Disconnect(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()

    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // Set Gin to release mode before creating the engine
    gin.SetMode(gin.ReleaseMode)

    // Initialize MongoDB client
    db := config.MongoClient.Database("sim_render")

    // Create default admin user if not exists
    usersCollection := db.Collection("users")
    var adminUser models.User
    err := usersCollection.FindOne(context.Background(), bson.M{"username": "admin69"}).Decode(&adminUser)
    if err == mongo.ErrNoDocuments {
        _, err := usersCollection.InsertOne(context.Background(), models.User{
            Username: "admin69",
            Password: "696969",
        })
        if err != nil {
            log.Printf("Error creating admin user: %v", err)
        }
    }

    // Initialize handler with MongoDB
    h := handlers.NewMongoHandler(db)

    // Initialize notification service with MongoDB
    notificationService := services.NewNotificationService(db)  // Changed from NewMongoNotificationService

    // Start notification checker in a goroutine
    go func() {
        for {
            if err := notificationService.CheckAndSendNotifications(); err != nil {
                fmt.Printf("Error sending notifications: %v\n", err)
            }
            time.Sleep(12 * time.Hour)
        }
    }()

    // Initialize router after setting release mode
    r := gin.Default()
    
    // Add CORS middleware
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

    // Serve static files first
    r.Static("/static", "./static")
    r.LoadHTMLGlob("static/*.html")

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        // Check MongoDB connection
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        err := db.Client().Ping(ctx, nil)
        if err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "Database ping failed"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    // Add check-user endpoint before protected routes
    r.GET("/api/check-user", func(c *gin.Context) {
        var user models.User
        err := db.Collection("users").FindOne(context.Background(), bson.M{"username": "admin69"}).Decode(&user)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "details": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "User exists", "username": user.Username})
    })
    
    // Public routes
    r.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })
    r.POST("/api/login", h.Login)

    // Protected routes
    protected := r.Group("/")
    protected.Use(authRequired)
    {
        protected.GET("/", func(c *gin.Context) {
            c.HTML(http.StatusOK, "index.html", nil)
        })
        
        api := protected.Group("/api")
        {
            api.POST("/sims", h.AddSim)
            api.GET("/sims", h.GetAllSims)
            api.PUT("/sims/:id", h.UpdateSimRechargeDate)
        }
    }

    // Start server
    port := config.GetPort()
    if port == "" {
        port = "8080" // fallback to default port
    }
    log.Printf("Server starting on port %s", port)
    r.Run("0.0.0.0:" + port)
}