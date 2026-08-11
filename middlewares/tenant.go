package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RequireOwnCompany ensures the {companyId} path parameter matches the
// caller's own company_id claim from their access token. Without this,
// any authenticated user could read or modify another company's data
// (and, combined with RequireAdminRole, another company's employees)
// simply by changing the ID in the URL.
func RequireOwnCompany(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlCompanyID := chi.URLParam(r, "companyId")
		tokenCompanyID, ok := r.Context().Value(CompanyIDKey).(string)

		if !ok || tokenCompanyID == "" || urlCompanyID != tokenCompanyID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
