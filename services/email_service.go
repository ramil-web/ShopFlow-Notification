package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(email, login, password string) error
}

type EmailService struct {
	username string
	password string
	host     string
	port     string
}

func NewEmailService(username, password, host, port string) *EmailService {
	return &EmailService{username: username, password: password, host: host, port: port}
}

func (s *EmailService) SendEmail(email, login, passwordStr string) error {
	to := []string{email}
	body := fmt.Sprintf("Subject: Registration\n\nHello %s,\nYour password is: %s", login, passwordStr)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// TLS connection
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.host,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer c.Quit()

	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err := c.Mail(s.username); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(body))
	if err != nil {
		return err
	}
	return w.Close()
}
