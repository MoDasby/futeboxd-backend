package email

import (
	"fmt"
	"net/smtp"
)

type EmailOpts struct {
	ContentType string
	To          string
	Subject     string
	Body        string
}

const (
	from string = "giullianomendes033@gmail.com"
	pass string = "kbmk itfr jzuq myny"
)

func SendMail(opts EmailOpts) error {
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", from, opts.To, opts.Subject, opts.ContentType, opts.Body)

	if err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("futeboxd@email.com", from, pass, "smtp.gmail.com"),
		from, []string{opts.To}, []byte(message)); err != nil {

		return err
	}

	return nil
}
