package db

import (
	"context"
	"testing"
	"time"
)

func TestRedisDatabaseGateway(t *testing.T) {
	ctx := context.Background()
	token := "2de0ZhksvEo6SItXrixCCqX54Gz9B0jhMUfkphlBrjy3cIYqaFdI33Z15uthyco2"
	limit := 100
	key := "127.0.0.1"
	ttl := time.Second
	t.Run("Lock", func(t *testing.T) {
		if err := redisConn.Lock(ctx, key, ttl); err != nil {
			t.Errorf("Lock(%v, %v, %v) = %v, want %v", ctx, key, ttl, err, nil)
		}
	})
	t.Run("IsLocked", func(t *testing.T) {
		if locked, err := redisConn.IsLocked(ctx, key); !locked || err != nil {
			t.Errorf("IsLocked(%v, %v) = (%v, %v), want (%v, %v)", ctx, key, locked, err, true, nil)
		}
		time.Sleep(ttl)
		if locked, err := redisConn.IsLocked(ctx, key); locked || err != nil {
			t.Errorf("IsLocked(%v, %v) = (%v, %v), want (%v, %v)", ctx, key, locked, err, false, nil)
		}
	})
	t.Run("SaveRequest", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if err := redisConn.SaveRequest(ctx, key); err != nil {
				t.Errorf("SaveRequest(%v, %v) = %v, want %v", ctx, key, err, nil)
			}
		}
	})
	t.Run("CountRequests", func(t *testing.T) {
		if total, err := redisConn.CountRequests(ctx, key); total != 10 || err != nil {
			t.Errorf("CountRequests(%v, %v) = (%v, %v), want (%v, %v)", ctx, key, total, err, 1, nil)
		}
	})
	t.Run("CreateToken", func(t *testing.T) {
		if err := redisConn.CreateToken(ctx, token, limit); err != nil {
			t.Errorf("CreateToken(%v, %v, %v) = %v, want %v", ctx, token, limit, err, nil)
		}
	})
	t.Run("GetTokenLimit", func(t *testing.T) {
		if got, err := redisConn.GetTokenLimit(ctx, token); got != limit || err != nil {
			t.Errorf("GetTokenLimit(%v, %v) = (%v, %v), want (%v, %v)", ctx, token, got, err, limit, nil)
		}
		token := "mTSdpLyrF0A2UpWSKmAk0I0VObEz3ocCdJBqIRU3HLHMrCYieJU61U4IvOB4xtrM"
		if got, err := redisConn.GetTokenLimit(ctx, token); got != 0 || err != nil {
			t.Errorf("GetTokenLimit(%v, %v) = (%v, %v), want (%v, %v)", ctx, token, got, err, 0, nil)
		}
	})
}
