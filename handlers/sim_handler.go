package handlers

import (
    "context"
    "net/http"
    "time"
    // "strings"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "mobilerecharge/models"
)

func (h *Handler) AddSim(c *gin.Context) {
    var sim models.Sim
    if err := c.ShouldBindJSON(&sim); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := h.db.Collection("sims").InsertOne(context.Background(), sim)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    sim.ID = result.InsertedID.(primitive.ObjectID)
    c.JSON(http.StatusOK, sim)
}

func (h *Handler) GetAllSims(c *gin.Context) {
    var sims []models.Sim
    ctx := context.Background()
    
    cursor, err := h.db.Collection("sims").Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    // Create a raw bson.D slice to store the documents
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Convert BSON documents to Sim structs with proper date handling
    for _, result := range results {
        sim := models.Sim{
            ID:       result["_id"].(primitive.ObjectID),
            Name:     result["name"].(string),
            Number:   result["number"].(string),
        }

        // Handle last_recharge_date
        if lrd, ok := result["last_recharge_date"].(string); ok {
            sim.LastRechargeDate = lrd
        }

        // Handle recharge_validity
        if rv, ok := result["recharge_validity"].(primitive.DateTime); ok {
            sim.RechargeValidity = rv.Time().Format(time.RFC3339)
        } else if rvStr, ok := result["recharge_validity"].(string); ok {
            sim.RechargeValidity = rvStr
        }

        // Handle incoming_call_validity
        if icv, ok := result["incoming_call_validity"].(primitive.DateTime); ok {
            sim.IncomingCallValidity = icv.Time().Format(time.RFC3339)
        } else if icvStr, ok := result["incoming_call_validity"].(string); ok {
            sim.IncomingCallValidity = icvStr
        }

        // Handle sim_expiry
        if se, ok := result["sim_expiry"].(primitive.DateTime); ok {
            sim.SimExpiry = se.Time().Format(time.RFC3339)
        } else if seStr, ok := result["sim_expiry"].(string); ok {
            sim.SimExpiry = seStr
        }

        sims = append(sims, sim)
    }

    c.JSON(http.StatusOK, sims)
}

func (h *Handler) UpdateSimRechargeDate(c *gin.Context) {
    id := c.Param("id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var updateData struct {
        LastRechargeDate     string `json:"last_recharge_date"`
        RechargeValidity     string `json:"recharge_validity"`
        IncomingCallValidity string `json:"incoming_call_validity"`
        SimExpiry            string `json:"sim_expiry"`
    }
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parse the last recharge date
    rechargeDate, err := time.Parse(time.RFC3339, updateData.LastRechargeDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
        return
    }

    // Calculate other dates
    update := bson.M{
        "$set": bson.M{
            "last_recharge_date":      rechargeDate.Format(time.RFC3339),
            "recharge_validity":       rechargeDate.AddDate(0, 0, 30).Format(time.RFC3339),
            "incoming_call_validity":  rechargeDate.AddDate(0, 0, 45).Format(time.RFC3339),
            "sim_expiry":             rechargeDate.AddDate(0, 0, 90).Format(time.RFC3339),
        },
    }

    result, err := h.db.Collection("sims").UpdateOne(
        context.Background(),
        bson.M{"_id": objectID},
        update,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "SIM not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "SIM updated successfully"})
}