package services

import (
    "context"

    "deployflow/shared/email"
)

type AuthEmailSender struct {
    sender email.Sender
}

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, email, token string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
}

func NewAuthEmailSender(sender email.Sender) EmailSender {
    return &AuthEmailSender{
        sender: sender,
    }
}

func (s *AuthEmailSender) SendVerificationEmail(ctx context.Context, emailAddress, token string) error {
    subject := "Verify your email"
    body := "Your email verification token is: " + token

    return s.sender.Send(ctx, emailAddress, subject, body)
}

func (s *AuthEmailSender) SendPasswordResetEmail(ctx context.Context, emailAddress, token string) error {
    subject := "Reset your password"
    body := "Your password reset token is: " + token

    return s.sender.Send(ctx, emailAddress, subject, body)
}