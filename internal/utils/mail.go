package utils

import (
	"fmt"
	"net/smtp"
	"os"
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
