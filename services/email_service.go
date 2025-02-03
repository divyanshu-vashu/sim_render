package services

import (
    "fmt"
    "net/smtp"
    "time"
)

const (
    smtpHost = "smtp.gmail.com"
    smtpPort = "587"
    fromEmail = "vdkalife@gmail.com"
    emailPassword = "xfmz rlod pixm mjvi"
)

var toEmails = []string{
    "vashusingh2005.jan@gmail.com",
    "divyanshusingh@appointy.com",
    "divyanshu.singh2021c@vitstudent.ac.in",
}

type EmailService struct {
    from     string
    password string
    host     string
    port     string
}

func NewEmailService() *EmailService {
    return &EmailService{
        from:     "vdkalife@gmail.com",  // Replace with your email
        password: "xfmz rlod pixm mjvi",   // Your app password
        host:     "smtp.gmail.com",
        port:     "587",
    }
}

func (s *EmailService) SendEmail(subject, message string) error {
    body := fmt.Sprintf("Subject: %s\n\n%s", subject, message)
    auth := smtp.PlainAuth("", s.from, s.password, s.host)
    addr := fmt.Sprintf("%s:%s", s.host, s.port)

    err := smtp.SendMail(addr, auth, s.from, toEmails, []byte(body))
    if err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }

    return nil
}

func (s *EmailService) SendExpiryNotification(simName string, number string, expiryDate time.Time, daysLeft int) error {
    auth := smtp.PlainAuth("", fromEmail, emailPassword, smtpHost)
    
    subject := "SIM Recharge Reminder"
    body := fmt.Sprintf(
        "Your SIM card %s (%s) will expire in %d days on %s. Please recharge soon to avoid service interruption.",
        simName,
        number,
        daysLeft,
        expiryDate.Format("2006-01-02"),
    )
    
    msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
    
    return smtp.SendMail(
        smtpHost+":"+smtpPort,
        auth,
        fromEmail,
        toEmails,
        msg,
    )
}