package email

import "context"

type Sender interface {
	Send(ctx context.Context, to, subject, body string) error
}