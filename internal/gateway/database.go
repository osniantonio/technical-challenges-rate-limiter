package gateway

import (
	"context"
	"time"
)

type DatabaseOptions struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type DatabaseGateway interface {
	Lock(context.Context, string, time.Duration) error
	IsLocked(context.Context, string) (bool, error)
	SaveRequest(context.Context, string) error
	CountRequests(context.Context, string) (int, error)
	CreateToken(context.Context, string, int) error
	GetTokenLimit(context.Context, string) (int, error)
}
