package repositories

import (
	"company_mgmt_api/internal/testutils"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateOTP(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewOTPRepository(db)

	mock.ExpectExec("INSERT INTO otps").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// FIX: Removed time.Now() from the arguments
	err := repo.Create(context.Background(), "user", "verify", "hash")
	if err != nil {
		t.Fatal(err)
	}
}
