package smtp

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

type SMTPConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	Password    string
	LogoURL     string
}

var DefaultConfig = &SMTPConfig{
	SMTPHost:    "smtp.gmail.com",
	SMTPPort:    "587",
	SenderEmail: "your-email@gmail.com",
	Password:    "your-app-password",
	LogoURL:     "https://app-logo-url",
}

type EmailData struct {
	Title       string
	Description string
	OTP         string
	Link        string
	LogoURL     string
}

func SendOTPEmail(recipient, otp string) error {

	Title := "Your OTP Code"
	Description := "Please use the following code to complete your login."

	emailData := EmailData{
		Title:       Title,
		Description: Description,
		OTP:         otp,
		LogoURL:     DefaultConfig.LogoURL,
	}

	emailBody, err := parseTemplate("pkg/smtp/email_template.html", emailData)
	if err != nil {
		return err
	}

	from := DefaultConfig.SenderEmail
	to := []string{recipient}
	subject := "Your OTP Code"
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", recipient, subject, emailBody))

	auth := smtp.PlainAuth("", DefaultConfig.SenderEmail, DefaultConfig.Password, DefaultConfig.SMTPHost)

	err = smtp.SendMail(DefaultConfig.SMTPHost+":"+DefaultConfig.SMTPPort, auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", recipient)
	return nil
}

func parseTemplate(templateFileName string, data EmailData) (string, error) {
	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return tpl.String(), nil
}
