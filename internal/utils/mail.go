package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/db"
)

func SendMail(toAddr string) error {

	from := "dev.test.jalen@gmail.com"
	to := "jalenarms@outlook.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		"This is a test"
	fmt.Println(os.Getenv("EMAIL_AP"))
	auth := smtp.PlainAuth("", "dev.test.jalen@gmail.com", os.Getenv("EMAIL_AP"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func SendOrderPaidEmail(order *db.Order) error {
	from := "dev.test.jalen@gmail.com"

	subject := "Thank you for your order!\n"
	contentType := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body><p>Test</p></body></html>"

	msg := []byte(subject + contentType + body)

	auth := smtp.PlainAuth("", "dev.test.jalen@gmail.com", os.Getenv("EMAIL_AP"), "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{*order.CustomerEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}
