package repositories

import (
	"company_mgmt_api/internal/testutils" // Ensure this path is correct
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestCreateEmployee(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewEmployeeRepository(db)

	// Create runs inside a transaction and inserts into both users and
	// profiles — the mock must expect the whole sequence, not just one exec.
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").
		WithArgs("id", "company", "email", "hash", "employee", sqlmock.AnyArg()).
		WillReturnResult(testutils.SuccessResult())
	mock.ExpectExec("INSERT INTO profiles").
		WithArgs("id", "first", "last", sqlmock.AnyArg()).
		WillReturnResult(testutils.SuccessResult())
	mock.ExpectCommit()

	err := repo.Create(
		context.Background(),
		"id", "company", "email",
		"hash", "first", "last", "employee",
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
