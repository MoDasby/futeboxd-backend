package email

import (
	"context"
)

type Metadata struct {
	To        string
	Subject   string
	Body      string
	LocalPart string
}

type Client interface {
	SendMail(context.Context, Metadata) error
}
