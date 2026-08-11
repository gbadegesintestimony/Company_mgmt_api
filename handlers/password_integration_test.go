package handlers

import (
	"company_mgmt_api/internal/testutils"
	"company_mgmt_api/repositories"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5" // Use chi/v5
)

func TestPasswordResetFlow(t *testing.T) {
	router, mock := setupPasswordTestServer(t)

	// No domain in the request body, so the handler looks up company
	// domain "" and should find nothing.
	mock.ExpectQuery("FROM companies").WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(
		"POST",
		"/v1/password/forgot",
		strings.NewReader(`{"email":"user@test.com"}`),
	)

	res := testutils.ExecuteRequest(router, req)

	// Ensure this matches your handler's StatusNoContent (204)
	if res.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", res.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// setupPasswordTestServer wires PasswordHandler to a sqlmock-backed
// CompanyRepository (a nil repo panics as soon as Forgot dereferences it).
func setupPasswordTestServer(t *testing.T) (chi.Router, sqlmock.Sqlmock) {
	t.Helper()
	db, mock := testutils.MockDB()
	t.Cleanup(func() { db.Close() })

	h := &PasswordHandler{
		CompanyRepo: repositories.NewCompanyRepository(db),
	}

	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/password/forgot", h.Forgot)
	})

	return r, mock
}
