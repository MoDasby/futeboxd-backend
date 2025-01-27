package email

import (
	"context"
	"fmt"
	"net/smtp"
)

const (
	from string = "giullianomendes033@gmail.com"
	pass string = "zlfy agrn sdad sxac"
)

type GmailClient struct {
	from string
	pass string
}

func NewGmailClient(from, pass string) Client {
	return &GmailClient{from: from, pass: pass}
}

func (gc *GmailClient) SendMail(ctx context.Context, metadata Metadata) error {
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", from, metadata.To, metadata.Subject, "text/html", metadata.Body)

	if err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("futeboxd@email.com", from, pass, "smtp.gmail.com"),
		from, []string{metadata.To}, []byte(message)); err != nil {

		return err
	}

	return nil
}
