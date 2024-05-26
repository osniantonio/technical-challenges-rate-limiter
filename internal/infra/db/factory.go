package db

import (
	"errors"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
)

type createDatabaseConnectionFn func(options *gateway.DatabaseOptions) (gateway.DatabaseGateway, error)

var (
	errUnknowDatabaseProtocol = errors.New("unknow database protocol")
	factories                 = make(map[string]createDatabaseConnectionFn)
)

func NewDatabaseConnection(options *gateway.DatabaseOptions) (gateway.DatabaseGateway, error) {
	if fn, ok := factories[options.Protocol]; ok {
		return fn(options)
	}
	return nil, errUnknowDatabaseProtocol
}

func init() {
	factories["redis"] = newRedisDatabaseGateway
}
