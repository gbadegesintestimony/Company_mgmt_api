package services

import (
	"company_mgmt_api/models"
	"company_mgmt_api/repositories"
	"company_mgmt_api/utils"
	"context"
	"errors"
	"time"
)

type OTPService struct {
	Repo  *repositories.OTPRepository
	Email *EmailService
}

func NewOTPService(repo *repositories.OTPRepository, email *EmailService) *OTPService {
	return &OTPService{
		Repo:  repo,
		Email: email,
	}
}

func (s *OTPService) GenerateAndSendOTP(
	ctx context.Context,
	userID, email, purpose string,
) error {
	if err := CanSendOTP(ctx, s.Repo, userID); err != nil {
		return err
	}

	otpCode := utils.GenerateOTP()
	hash := utils.HashOTP(otpCode)

	otp := &models.OTP{
		UserID:    userID,
		CodeHash:  hash,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.Repo.Create(ctx, otp); err != nil {
		return err
	}

	if err := s.Email.SendOTP(ctx, email, otpCode, purpose); err != nil {
		return err
	}
	return nil
}

func (s *OTPService) VerifyOTP(
	ctx context.Context,
	userID, code, purpose string,
) error {
	otp, err := s.Repo.GetActive(ctx, userID, purpose)
	if err != nil {
		return errors.New("invalid code")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("code has expired")
	}

	if err := utils.VerifyOTP(otp.CodeHash, code); err != nil {
		return errors.New("invalid code")
	}

	return s.Repo.Consume(ctx, userID, purpose)
}
