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
	userID, rawCode, newPassword string,
) error {

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// hash := utils.HashOTP(rawCode)

	if err := s.OTP.VerifyAndConsumeTx(ctx, tx, userID, "password_reset", rawCode); err != nil {
		return err
	}

	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE users SET password_hash = $1 WHERE id = $2
	`, newHash, userID); err != nil {
		return err
	}

	if err := s.Sessions.RevokeAll(ctx, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func NewAuthService(db *sql.DB, otpRepo *repositories.OTPRepository, sessionRepo *repositories.SessionRepository) *AuthService {
	return &AuthService{
		DB:       db,
		OTP:      otpRepo,
		Sessions: sessionRepo,
	}
}
