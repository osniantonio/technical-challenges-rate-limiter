package ratelimiter

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
)

type DefaultRateLimiter struct {
	settings *Settings
	db       gateway.DatabaseGateway
}

func NewDefaultRateLimiter(settings *Settings, db gateway.DatabaseGateway) *DefaultRateLimiter {
	return &DefaultRateLimiter{settings, db}
}

func (rt *DefaultRateLimiter) CanGo(ctx context.Context, r *http.Request) (bool, error) {
	key := r.Header.Get("API_KEY")
	if key == "" || !rt.settings.LimitByToken {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false, err
		}
		key = host
	}
	locked, err := rt.db.IsLocked(ctx, key)
	if locked {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	limit, err := rt.db.GetTokenLimit(ctx, key)
	if err != nil {
		return false, err
	}
	if limit == 0 {
		limit = rt.settings.Ratelimit
	}
	total, err := rt.db.CountRequests(ctx, key)
	if err != nil {
		return false, err
	}
	err = rt.db.SaveRequest(ctx, key)
	if err != nil {
		return false, err
	}
	if total >= limit {
		if err := rt.db.Lock(ctx, key, time.Second*time.Duration(rt.settings.ExpirationTime)); err != nil {
			return false, err
		}
		return false, nil
	}
	return total <= limit, nil
}
