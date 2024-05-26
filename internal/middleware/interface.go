package middleware

import (
	"context"
	"net/http"
)

type Middleware interface {
	Execute(context.Context, http.Handler) http.Handler
}
