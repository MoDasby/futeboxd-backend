package email

import (
	"context"
	"fmt"

	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/resend/resend-go/v2"
)

type ResendClient struct {
	Client resend.Client
}

func NewResendClient(cfg config.Resend) Client {
	return &ResendClient{Client: *resend.NewClient(cfg.ApiKey)}
}

func (rc *ResendClient) SendMail(ctx context.Context, metadata Metadata) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s@mail.futeboxd.com", metadata.LocalPart),
		To:      []string{metadata.To},
		Subject: metadata.Subject,
		Html:    metadata.Body,
	}

	if _, err := rc.Client.Emails.SendWithContext(ctx, params); err != nil {
		return err
	}

	return nil
}
