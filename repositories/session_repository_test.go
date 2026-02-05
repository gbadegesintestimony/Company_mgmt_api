package repositories

import (
	"company_mgmt_api/internal/testutils"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSessionCreate(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewSessionRepository(db)

	mock.ExpectExec("INSERT INTO sessions").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// FIX: Changed repo.Create to repo.CreateSession
	err := repo.CreateSession(context.Background(), "user", "hash")
	if err != nil {
		t.Fatal(err)
	}
}
