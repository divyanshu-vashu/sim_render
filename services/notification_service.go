package services

import (
    "fmt"
    "time"
    "gorm.io/gorm"
    "mobilerecharge/models"
)

type NotificationService struct {
    db *gorm.DB
    emailService *EmailService
}

func NewNotificationService(db *gorm.DB) *NotificationService {
    return &NotificationService{
        db: db,
        emailService: NewEmailService(),
    }
}

func (s *NotificationService) CheckAndSendNotifications() error {
    var sims []models.Sim
    if err := s.db.Find(&sims).Error; err != nil {
        return fmt.Errorf("error fetching sims: %v", err)
    }

    for _, sim := range sims {
        // Check if recharge validity date is set
        if !sim.RechargeValidity.IsZero() {
            daysLeft := int(time.Until(sim.RechargeValidity).Hours() / 24)
            if daysLeft <= 7 && daysLeft >= 0 {
                err := s.emailService.SendExpiryNotification(
                    sim.Name,
                    sim.Number,
                    sim.RechargeValidity,
                    daysLeft,
                )
                if err != nil {
                    fmt.Printf("Error sending notification for sim %s: %v\n", sim.Name, err)
                }
            }
        }
    }
    return nil
}