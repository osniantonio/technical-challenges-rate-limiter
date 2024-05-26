package ratelimiter

import "testing"

func TestNewSettings(t *testing.T) {
	rateLimit := 10
	expirationTime := 60
	limitByToken := true
	settings := NewSettings(rateLimit, expirationTime, limitByToken)
	if settings.Ratelimit != rateLimit {
		t.Errorf("settings.rateLimit = %v, want %v", settings.Ratelimit, rateLimit)
	}
	if settings.ExpirationTime != expirationTime {
		t.Errorf("settings.expirationTime = %v, want %v", settings.ExpirationTime, expirationTime)
	}
	if settings.LimitByToken != limitByToken {
		t.Errorf("settings.limitByToken = %v, want %v", settings.LimitByToken, limitByToken)
	}
}
