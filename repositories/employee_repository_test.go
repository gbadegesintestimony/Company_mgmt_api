package repositories

import (
	"company_mgmt_api/internal/testutils" // Ensure this path is correct
	"context"
	"testing"
)

func TestCreateEmployee(t *testing.T) {
	db, mock := testutils.MockDB()
	defer db.Close()

	repo := NewEmployeeRepository(db)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("id", "company", "email", "hash", "first", "last", "employee").
		WillReturnResult(testutils.SuccessResult())

	err := repo.Create(
		context.Background(),
		"id", "company", "email",
		"hash", "first", "last", "employee",
	)

	if err != nil {
		t.Fatal(err)
	}
}
