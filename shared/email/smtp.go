package email

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPSender(host string, port int, username, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	message := []byte(
		"From: " + s.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body,
	)

	address := net.JoinHostPort(s.host, strconv.Itoa(s.port))

	if err := smtp.SendMail(address, auth, s.from, []string{to}, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

var _ Sender = (*SMTPSender)(nil)