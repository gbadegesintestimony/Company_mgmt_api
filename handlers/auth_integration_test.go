package handlers

import (
	"company_mgmt_api/config"
	"company_mgmt_api/internal/testutils"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestAdminRegister(t *testing.T) {
	router, mock := setupAuthTestServer(t)

	// No existing company for domain "" (none provided in the request body).
	mock.ExpectQuery("FROM companies").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO companies").WillReturnResult(sqlmock.NewResult(1, 1))

	// CreateAdmin runs in a transaction: users row, then profiles row.
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO profiles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// OTP generation: rate-limit check, then insert.
	mock.ExpectQuery("FROM otps").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO otps").WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(
		"POST",
		"/v1/auth/admin/register",
		strings.NewReader(`{"email":"admin@test.com","password":"secret123"}`),
	)

	res := testutils.ExecuteRequest(router, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", res.Code, res.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// setupAuthTestServer wires a real AuthHandler to a sqlmock-backed DB and a
// dev-mode (log-only) email sender, so Register can be exercised without a
// live Postgres instance or a real Resend API key.
func setupAuthTestServer(t *testing.T) (chi.Router, sqlmock.Sqlmock) {
	t.Helper()
	db, mock := testutils.MockDB()
	t.Cleanup(func() { db.Close() })

	emailSvc := services.NewDevEmailService("test@example.com")
	otpRepo := repositories.NewOTPRepository(db)

	h := &AuthHandler{
		Cfg:         &config.Config{},
		Email:       emailSvc,
		Sessions:    repositories.NewSessionRepository(db),
		UserRepo:    repositories.NewUserRepository(db),
		CompanyRepo: repositories.NewCompanyRepository(db),
		OTPService:  services.NewOTPService(otpRepo, emailSvc),
	}

	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth/admin/register", h.Register)
	})

	return r, mock
}
