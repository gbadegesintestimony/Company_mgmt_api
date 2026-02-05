package repositories

import (
	"company_mgmt_api/internal/testutils"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuditLog(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewAuditRepository(db)

	mock.ExpectExec("INSERT INTO audit_logs").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// FIX: Changed 'nil' to '"{}"' because the function expects a string
	err := repo.Log(context.Background(), "ACTION", "actor", "target", "{}")
	if err != nil {
		t.Fatal(err)
	}
}
