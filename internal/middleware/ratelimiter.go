package middleware

import (
	"context"
	"net/http"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/ratelimiter"
)

var (
	tooManyRequestsMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

type RateLimiterMiddleware struct {
	rt ratelimiter.RateLimiter
}

func NewRateLimiterMiddleware(rt ratelimiter.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{rt}
}

func (m *RateLimiterMiddleware) Execute(ctx context.Context, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		proceed, err := m.rt.CanGo(ctx, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !proceed {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(tooManyRequestsMessage))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
