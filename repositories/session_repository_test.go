package repositories

import (
	"company_mgmt_api/internal/testutils"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSessionCreate(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewSessionRepository(db)

	mock.ExpectExec("INSERT INTO sessions").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateSession(
		context.Background(),
		"user-123",
		"company-123",
		"admin",
		"refresh-token-hash",
		time.Now().Add(24*time.Hour),
	)
	if err != nil {
		t.Fatal(err)
	}
}
