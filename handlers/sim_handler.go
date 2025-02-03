package handlers

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "mobilerecharge/models"
)

type Handler struct {
    DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{DB: db}
}

func (h *Handler) AddSim(c *gin.Context) {
    var sim models.Sim
    if err := c.ShouldBindJSON(&sim); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set default dates if not provided
    now := time.Now()
    if sim.LastRechargeDate.IsZero() {
        sim.LastRechargeDate = now
    }
    if sim.RechargeValidity.IsZero() {
        sim.RechargeValidity = now.Add(30 * 24 * time.Hour) // 30 days validity
    }
    if sim.IncomingCallValidity.IsZero() {
        sim.IncomingCallValidity = now.Add(45 * 24 * time.Hour) // 45 days validity
    }
    if sim.SimExpiry.IsZero() {
        sim.SimExpiry = now.Add(90 * 24 * time.Hour) // 90 days validity
    }

    // Debug log
    fmt.Printf("Adding SIM with data: %+v\n", sim)

    if err := h.DB.Create(&sim).Error; err != nil {
        fmt.Printf("Error creating SIM: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, sim)
}

func (h *Handler) GetAllSims(c *gin.Context) {
    var sims []models.Sim
    result := h.DB.Order("last_recharge_date desc").Find(&sims)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, sims)
}

// Add this new handler function
func (h *Handler) UpdateSimRechargeDate(c *gin.Context) {
    var updateData struct {
        LastRechargeDate     string `json:"last_recharge_date"`
        RechargeValidity     string `json:"recharge_validity"`
        IncomingCallValidity string `json:"incoming_call_validity"`
        SimExpiry            string `json:"sim_expiry"`
    }

    if err := c.ShouldBindJSON(&updateData); err != nil {
        fmt.Printf("JSON binding error: %v\n", err)
        fmt.Printf("Received data: %+v\n", updateData)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    id := c.Param("id")
    var sim models.Sim
    if err := h.DB.First(&sim, id).Error; err != nil {
        fmt.Printf("Database fetch error for ID %s: %v\n", id, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Sim not found"})
        return
    }

    // Parse and validate dates with error logging
    layout := "2006-01-02"
    
    if lastRecharge, err := time.Parse(layout, updateData.LastRechargeDate); err == nil {
        sim.LastRechargeDate = lastRecharge
    } else {
        fmt.Printf("Error parsing LastRechargeDate: %v\n", err)
    }
    
    if rechargeVal, err := time.Parse(layout, updateData.RechargeValidity); err == nil {
        sim.RechargeValidity = rechargeVal
    } else {
        fmt.Printf("Error parsing RechargeValidity: %v\n", err)
    }
    
    if incomingVal, err := time.Parse(layout, updateData.IncomingCallValidity); err == nil {
        sim.IncomingCallValidity = incomingVal
    } else {
        fmt.Printf("Error parsing IncomingCallValidity: %v\n", err)
    }
    
    if simExp, err := time.Parse(layout, updateData.SimExpiry); err == nil {
        sim.SimExpiry = simExp
    } else {
        fmt.Printf("Error parsing SimExpiry: %v\n", err)
    }

    fmt.Printf("Updating SIM with data: %+v\n", sim)
    if err := h.DB.Save(&sim).Error; err != nil {
        fmt.Printf("Database save error: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, sim)
}