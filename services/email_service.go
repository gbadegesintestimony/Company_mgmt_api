package services

import (
	"context"
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client  *resend.Client
	from    string
	devMode bool
}

func NewEmailService(apiKey, from string) *EmailService {
	client := resend.NewClient(apiKey)
	return &EmailService{
		client: client,
		from:   from,
	}
}

// NewDevEmailService returns an EmailService that logs the OTP instead of
// sending real email, for local development and tests where no live email
// provider credentials are available.
func NewDevEmailService(from string) *EmailService {
	return &EmailService{from: from, devMode: true}
}

// SendOTP sends an OTP email for verification or password reset
func (s *EmailService) SendOTP(ctx context.Context, to, otp, purpose string) error {
	if s.devMode {
		log.Printf("DEV MODE OTP — to: %s | purpose: %s | code: %s", to, purpose, otp)
		return nil
	}
	subject := fmt.Sprintf("%s Verification Code", purpose)

	body := fmt.Sprintf(`
	<h2> %s Verification</h2>
	<p>Your verificationcode is:</p>
	<h1>%s</h1>
	<p>This code expires in 10 minutes.</p>
	`, purpose, otp)

	_, err := s.client.Emails.SendWithContext(
		ctx,
		&resend.SendEmailRequest{
			From:    s.from,
			To:      []string{to},
			Subject: subject,
			Html:    body,
		},
	)

	return err
}
