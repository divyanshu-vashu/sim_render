package services

import (
    "context"
    "time"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "mobilerecharge/models"
)

type NotificationService struct {
    db *mongo.Database
    emailService *EmailService
}

func NewNotificationService(db *mongo.Database) *NotificationService {
    return &NotificationService{
        db: db,
        emailService: NewEmailService(),
    }
}

func (s *NotificationService) CheckAndSendNotifications() error {
    var sims []models.Sim
    cursor, err := s.db.Collection("sims").Find(context.Background(), bson.M{})
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())

    if err = cursor.All(context.Background(), &sims); err != nil {
        return err
    }

    now := time.Now()
    for _, sim := range sims {
        // Parse the sim expiry date
        simExpiry, err := time.Parse(time.RFC3339, sim.SimExpiry)
        if err != nil {
            continue // Skip this sim if date parsing fails
        }

        // Check for 2 days before expiry
        twoDaysBeforeExpiry := simExpiry.AddDate(0, 0, -2)
        if now.After(twoDaysBeforeExpiry) && now.Before(simExpiry) {
            message := fmt.Sprintf("SIM %s (%s) will expire in 2 days on %s", 
                sim.Name, sim.Number, simExpiry.Format("2006-01-02"))
            if err := s.emailService.SendEmail("SIM Expiry Alert - 2 Days", message); err != nil {
                fmt.Printf("Failed to send 2-day notification for SIM %s: %v\n", sim.Number, err)
            }
        }

        // Check for 1 day before expiry
        oneDayBeforeExpiry := simExpiry.AddDate(0, 0, -1)
        if now.After(oneDayBeforeExpiry) && now.Before(simExpiry) {
            message := fmt.Sprintf("URGENT: SIM %s (%s) will expire tomorrow on %s", 
                sim.Name, sim.Number, simExpiry.Format("2006-01-02"))
            if err := s.emailService.SendEmail("URGENT: SIM Expiry Alert - 1 Day", message); err != nil {
                fmt.Printf("Failed to send 1-day notification for SIM %s: %v\n", sim.Number, err)
            }
        }
    }

    return nil
}