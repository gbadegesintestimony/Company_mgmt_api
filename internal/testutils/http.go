package testutils

import (
	"net/http"
	"net/http/httptest"
)

// ExecuteRequest remains the same
func ExecuteRequest(handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
