package repositories

import (
	"company_mgmt_api/internal/testutils"
	"company_mgmt_api/models"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateOTP(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewOTPRepository(db)

	mock.ExpectExec("INSERT INTO otps").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create an OTP model and pass to repository
	otp := &models.OTP{
		UserID:    "user",
		Purpose:   "verify",
		CodeHash:  "hash",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	err := repo.Create(context.Background(), otp)
	if err != nil {
		t.Fatal(err)
	}
}
