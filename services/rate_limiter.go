package services

import (
	"context"
	"fmt"
	"time"

	"company_mgmt_api/repositories"
)

func CanSendOTP(
	ctx context.Context,
	repo *repositories.OTPRepository,
	userID string,
) error {

	count, err := repo.CountRecent(ctx, userID, 15*time.Minute)
	if err != nil {
		return err
	}
	if count >= 5 {
		return fmt.Errorf("too many OTP requests, please try again later")
	}
	return nil
}
