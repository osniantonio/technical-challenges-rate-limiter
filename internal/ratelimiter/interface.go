package ratelimiter

import (
	"context"
	"net/http"
)

type RateLimiter interface {
	CanGo(context.Context, *http.Request) (bool, error)
}
