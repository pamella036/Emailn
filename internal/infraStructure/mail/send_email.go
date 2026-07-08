package mail

import (
	"emailn/internal/domain/campaign"
	"fmt"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

func SendMail(campaign *campaign.Campaign) error {
	fmt.Println("Send mail...")

	start := time.Now()
	d := gomail.NewDialer("EMAIL_SMTP", 587, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))
	duration := time.Since(start)
	fmt.Println("Dialer created in", duration)

	var emails []string
	for _, contact := range campaign.Contacts {
		emails = append(emails, contact.Email)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", emails...)
	m.SetHeader("Subject", campaign.Name)
	m.SetBody("text/html", campaign.Content)

	start = time.Now()
	err := d.DialAndSend(m)
	duration = time.Since(start)
	fmt.Println("Dialer created in", duration)

	return err
}
