package db

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	redisOptions *gateway.DatabaseOptions
	redisConn    gateway.DatabaseGateway
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	redisOptions = &gateway.DatabaseOptions{
		Protocol: "redis",
		Host:     host,
		Port:     port,
		Database: "0",
	}
	os.Exit(m.Run())
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func TestNewDatabaseConnection(t *testing.T) {
	var err error
	if redisConn, err = NewDatabaseConnection(redisOptions); err != nil {
		log.Fatal(err)
	}
	options := &gateway.DatabaseOptions{
		Protocol: "mysql",
		Host:     "localhost",
		Port:     "3306",
		User:     "",
		Password: "",
		Database: "",
	}
	if _, err := NewDatabaseConnection(options); err == nil {
		t.Errorf("NewDatabaseConnection(%v) = %v, want %v", options, err, errUnknowDatabaseProtocol)
	}
}
