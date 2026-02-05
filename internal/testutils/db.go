package testutils

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

// Add this function:
func SuccessResult() sql.Result {
	return sqlmock.NewResult(1, 1)
}
