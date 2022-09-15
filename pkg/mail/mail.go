package mail

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

// EmailInput represents an email input
type EmailInput struct {
	To         []string
	Subject    string
	Message    string
	Attachment []string
}

// SendEmail sends an email with provided parameter
func SendEmail(input EmailInput) error {
	// Get mail configuration from environment variables
	mailUsername := os.Getenv("MAIL_USERNAME")
	if _, err := mail.ParseAddress(mailUsername); err != nil {
		return fmt.Errorf("email address is invalid: %v", err)
	}

	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort := os.Getenv("MAIL_PORT")

	// Validate input
	for _, email := range input.To {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf("email address is invalid: %v", err)
		}
	}
	if input.Subject == "" {
		return fmt.Errorf("subject cannot be empty")
	}
	if input.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Setup information to email
	mail := email.NewEmail()
	mail.From = mailUsername
	mail.To = input.To
	mail.Subject = input.Subject
	mail.HTML = []byte(input.Message)

	// Attach file from list file
	for _, filePath := range input.Attachment {
		if _, err := mail.AttachFile(filePath); err != nil {
			return fmt.Errorf("cannot attach file: %v", err)
		}
	}

	// Authenticate and send email
	auth := smtp.PlainAuth("", mailUsername, mailPassword, mailHost)
	if err := mail.Send(mailHost+":"+mailPort, auth); err != nil {
		return fmt.Errorf("cannot send email: %v", err)
	}
	return nil
}
