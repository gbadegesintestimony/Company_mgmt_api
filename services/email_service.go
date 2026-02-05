package services

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	// Fields for email service configuration
	client *resend.Client
	from   string
}

func NewEmailService(apiKey, from string) *EmailService {
	client := resend.NewClient(apiKey)
	return &EmailService{
		client: client,
		from:   from,
	}
}

// SendOTP sends an OTP emaiil for verification or password reset

func (s *EmailService) SendOTP(ctx context.Context, to, otp, purpose string) error {
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
