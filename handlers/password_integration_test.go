package handlers

import (
	"company_mgmt_api/internal/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5" // Use chi/v5
)

func TestPasswordResetFlow(t *testing.T) {
	router := setupTestServer()

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
}

// Helper function updated for Chi
func setupTestServer() chi.Router {
	r := chi.NewRouter()

	h := &PasswordHandler{
		// OTP: mockService, (You'll need to inject a real or mock service here)
	}

	r.Route("/v1", func(r chi.Router) {
		r.Post("/password/forgot", h.Forgot)
	})

	return r
}
