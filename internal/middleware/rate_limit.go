package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bjdms/api/pkg/response"
)

// RateLimiter implements a simple in-memory leaky bucket rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Limit middleware restricts requests by IP
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Simplistic, in production use X-Forwarded-For

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Clean up old requests
		validRequests := []time.Time{}
		for _, reqTime := range rl.requests[ip] {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) >= rl.limit {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
			response.Error(w, http.StatusTooManyRequests, "too_many_requests", "Rate limit exceeded. Please try again later.", "")
			return
		}

		rl.requests[ip] = append(validRequests, now)
		next.ServeHTTP(w, r)
	})
}
