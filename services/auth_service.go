package services

import (
	"context"
	"database/sql"

	"company_mgmt_api/repositories"
	"company_mgmt_api/utils"
)

type AuthService struct {
	DB       *sql.DB
	OTP      *repositories.OTPRepository
	Sessions *repositories.SessionRepository
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	userID, otpHash, newPassword string,
) error {

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.OTP.VerifyAndConsumeTx(ctx, tx, userID, "password_reset", otpHash); err != nil {
		return err
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE users SET password_hash = $1 WHERE id = $2
	`, hash, userID); err != nil {
		return err
	}

	if err := s.Sessions.RevokeAll(ctx, userID); err != nil {
		return err
	}

	return tx.Commit()
}
