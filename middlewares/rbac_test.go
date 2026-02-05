package middlewares

import (
	"context" // Add this import
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminOnlyMiddleware(t *testing.T) {
	// 1. Match the function name in rbac.go
	handler := RequireAdminRole(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)

	// 2. Use the helper function defined below
	req = req.WithContext(contextWithRole(req.Context(), "employee"))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden, got %d", rr.Code)
	}
}

// Helper to mock the context role
func contextWithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}
