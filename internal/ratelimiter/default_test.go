package ratelimiter

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/infra/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbOptions *gateway.DatabaseOptions
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
	dbOptions = &gateway.DatabaseOptions{
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

func TestDefaultRateLimiter(t *testing.T) {
	ctx := context.Background()
	token := "g8dXuf2MqNkqJ5tb47qw4m6thqYbrsK24SFZV4OiS83Lmbp8NCYulXtO3tyHJyZN"
	limit := 100
	settings := NewSettings(10, 60, true)
	conn, err := db.NewDatabaseConnection(dbOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.CreateToken(ctx, token, limit); err != nil {
		t.Fatal(err)
	}
	rt := NewDefaultRateLimiter(settings, conn)
	t.Run("Execute", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:8080"
		for i := 0; i < 10; i++ {
			if proceed, err := rt.Execute(ctx, req); !proceed || err != nil {
				t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, true, nil)
			}
		}
		if proceed, err := rt.Execute(ctx, req); proceed || err != nil {
			t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
		}
		req.Header.Add("API_KEY", token)
		for i := 0; i < limit; i++ {
			if proceed, err := rt.Execute(ctx, req); !proceed || err != nil {
				t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, true, nil)
			}
		}
		if proceed, err := rt.Execute(ctx, req); proceed || err != nil {
			t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
		}
	})
}
