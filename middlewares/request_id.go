package middlewares

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	// Implementation for request ID middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate or extract request ID and set it in the context or headers
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)

	})
}
