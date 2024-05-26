package ratelimiter

import (
	"context"
	"net/http"
)

type RateLimiter interface {
	Execute(context.Context, *http.Request) (bool, error)
}
