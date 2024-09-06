package models

import (
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type Email struct {
	Id      uuid.UUID `json:"id" db:"Id"`
	SentTo  []string  `json:"sentTo" db:"SentTo"`
	Subject string
	Body    string
	SentAt  time.Time
}

func NewEmail(sendTo []string, subject, body string) *Email {
	newId := generateUUIDv4()

	return &Email{
		Id:      newId,
		SentTo:  sendTo,
		Subject: subject,
		Body:    body,
		SentAt:  time.Now(),
	}
}

func (e *Email) SendMail() error {
	from := "dev.test.jalen@gmail.com"
	password := os.Getenv("EMAIL_PASS")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	m := []byte(fmt.Sprintf("Subject: %s\n\n%s", e.Subject, e.Body))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, e.SentTo, m)
	if err != nil {
		return err
	}

	return nil
}
