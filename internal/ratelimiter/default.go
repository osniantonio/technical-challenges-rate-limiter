package ratelimiter

import (
	"context"
	"log"
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

func (rt *DefaultRateLimiter) ExecuteByToken(ctx context.Context, r *http.Request) (bool, error) {
	key := r.Header.Get("API_KEY")
	log.Printf("Executing rate limit by token for key: %s", key)
	if key == "" || !rt.settings.LimitByToken {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error splitting host and port: %v", err)
			return false, err
		}
		key = host
	}
	locked, err := rt.db.IsLocked(ctx, key)
	if locked {
		log.Printf("Key %s is locked", key)
		return false, nil
	}
	if err != nil {
		log.Printf("Error checking if key %s is locked: %v", key, err)
		return false, err
	}
	limit, err := rt.db.GetTokenLimit(ctx, key)
	if err != nil {
		log.Printf("Error getting token limit for key %s: %v", key, err)
		return false, err
	}
	if limit == 0 {
		limit = rt.settings.Ratelimit
	}
	total, err := rt.db.CountRequests(ctx, key)
	if err != nil {
		log.Printf("Error counting requests for key %s: %v", key, err)
		return false, err
	}
	err = rt.db.SaveRequest(ctx, key)
	if err != nil {
		log.Printf("Error saving request for key %s: %v", key, err)
		return false, err
	}
	if total >= limit {
		log.Printf("Rate limit exceeded for key %s. Total: %d, Limit: %d", key, total, limit)
		if err := rt.db.Lock(ctx, key, time.Second*time.Duration(rt.settings.ExpirationTime)); err != nil {
			log.Printf("Error locking key %s: %v", key, err)
			return false, err
		}
		return false, nil
	}
	log.Printf("Request allowed for key %s. Total: %d, Limit: %d", key, total, limit)
	return total <= limit, nil
}

func (rt *DefaultRateLimiter) ExecuteByIP(ctx context.Context, r *http.Request) (bool, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error splitting host and port: %v", err)
		return false, err
	}
	log.Printf("Executing rate limit by IP for host: %s", host)

	locked, err := rt.db.IsLocked(ctx, host)
	if err != nil {
		log.Printf("Error checking if host is locked: %v", err)
		return false, err
	}
	if locked {
		log.Printf("Host %s is locked", host)
		return false, nil
	}

	limit := rt.settings.Ratelimit

	total, err := rt.db.CountRequests(ctx, host)
	if err != nil {
		log.Printf("Error counting requests for host %s: %v", host, err)
		return false, err
	}

	if total >= limit {
		log.Printf("Rate limit exceeded for host %s. Total: %d, Limit: %d", host, total, limit)
		err := rt.db.Lock(ctx, host, time.Second*time.Duration(rt.settings.ExpirationTime))
		if err != nil {
			log.Printf("Error locking host %s: %v", host, err)
			return false, err
		}
		return false, nil
	}

	log.Printf("Request allowed for host %s. Total: %d, Limit: %d", host, total, limit)
	err = rt.db.SaveRequest(ctx, host)
	if err != nil {
		log.Printf("Error saving request for host %s: %v", host, err)
		return false, err
	}

	return true, nil
}

func (rt *DefaultRateLimiter) isTokenLimit(r *http.Request) bool {
	isToken := rt.settings.LimitByToken && r.Header.Get("API_KEY") != ""
	log.Printf("Is token limit: %v", isToken)
	return isToken
}

func (rt *DefaultRateLimiter) Execute(ctx context.Context, r *http.Request) (bool, error) {
	if rt.isTokenLimit(r) {
		return rt.ExecuteByToken(ctx, r)
	} else {
		return rt.ExecuteByIP(ctx, r)
	}
}
