package middlewares

import (
	"company_mgmt_api/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	RoleKey      contextKey = "role"
	CompanyIDKey contextKey = "company_id"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implementation for authentication middleware
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := utils.ParseToken(auth[7:], secret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims["user_id"])
			ctx = context.WithValue(ctx, RoleKey, claims["role"])
			ctx = context.WithValue(ctx, CompanyIDKey, claims["company_id"])

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
