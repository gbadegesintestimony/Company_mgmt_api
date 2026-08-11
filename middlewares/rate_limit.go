package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type loginAttempts struct {
	count     int
	windowEnd time.Time
}

// LoginRateLimiter throttles requests per client IP to `limit` within
// `window`, returning 429 once exceeded. It guards login endpoints, which
// have no other brute-force protection.
//
// This is in-memory and per-instance: fine for a single API replica, but a
// multi-instance deployment would need a shared store (e.g. Redis) for the
// limit to hold across instances. Place it behind middleware.RealIP so
// r.RemoteAddr reflects the real client IP.
func LoginRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	attempts := make(map[string]*loginAttempts)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			now := time.Now()

			mu.Lock()
			a, ok := attempts[ip]
			if !ok || now.After(a.windowEnd) {
				a = &loginAttempts{windowEnd: now.Add(window)}
				attempts[ip] = a
			}
			a.count++
			blocked := a.count > limit
			mu.Unlock()

			if blocked {
				http.Error(w, "too many attempts, please try again later", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
