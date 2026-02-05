package handlers

import (
	"company_mgmt_api/internal/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdminRegister(t *testing.T) {
	router := setupTestServer()

	req := httptest.NewRequest(
		"POST",
		"/v1/auth/admin/register",
		strings.NewReader(`{"email":"admin@test.com","password":"secret"}`),
	)

	res := testutils.ExecuteRequest(router, req)
	if res.Code != http.StatusCreated {
		t.Fatal("expected 201")
	}
}
