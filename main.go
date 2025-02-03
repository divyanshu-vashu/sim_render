package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "gorm.io/gorm/logger"
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
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // Set Gin to release mode before creating the engine
    gin.SetMode(gin.ReleaseMode)

    // Initialize database
    dsn := config.GetDBConfig()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        panic("failed to connect database: " + err.Error())
    }
    
    // Migrate the schema
    db.AutoMigrate(&models.Sim{}, &models.User{})  // Combine migrations
    
    // Create default admin user if not exists
    var adminUser models.User
    if db.Where("username = ?", "admin69").First(&adminUser).Error != nil {
        db.Create(&models.User{
            Username: "admin69",
            Password: "696969",
        })
    }

    // Initialize handler
    h := handlers.NewHandler(db)

    // Initialize notification service
    notificationService := services.NewNotificationService(db)

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
        // Check database connection
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "Database connection error"})
            return
        }
        
        // Ping database
        if err := sqlDB.Ping(); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "Database ping failed"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    // Add check-user endpoint before protected routes
    r.GET("/api/check-user", func(c *gin.Context) {
        var user models.User
        result := db.Where("username = ?", "admin69").First(&user)
        if result.Error != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "details": result.Error.Error()})
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