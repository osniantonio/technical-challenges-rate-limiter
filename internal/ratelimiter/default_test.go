package ratelimiter_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/infra/db"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/ratelimiter"
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

	code := m.Run()
	os.Exit(code)

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func TestExecutionTheRateLimiter(t *testing.T) {
	ctx := context.Background()
	token := "g8dXuf2MqNkqJ5tb47qw4m6thqYbrsK24SFZV4OiS83Lmbp8NCYulXtO3tyHJyZN"
	limit := 100
	settings := ratelimiter.NewSettings(10, 60, true)
	conn, err := db.NewDatabaseConnection(dbOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.CreateToken(ctx, token, limit); err != nil {
		t.Fatal(err)
	}
	rt := ratelimiter.NewDefaultRateLimiter(settings, conn)
	t.Run("Execute by IP", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:8080"

		for i := 0; i < settings.Ratelimit; i++ {
			if proceed, err := rt.Execute(ctx, req); !proceed || err != nil {
				t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, true, nil)
			}
		}

		if proceed, err := rt.Execute(ctx, req); proceed || err != nil {
			t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
		}
	})

	t.Run("Execute by Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
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

func TestExceedingTheRateLimiter(t *testing.T) {
	ctx := context.Background()
	token := "g8dXuf2MqNkqJ5tb47qw4m6thqYbrsK24SFZV4OiS83Lmbp8NCYulXtO3tyHJyZN"
	limit := 100
	settings := ratelimiter.NewSettings(10, 60, true)
	conn, err := db.NewDatabaseConnection(dbOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.CreateToken(ctx, token, limit); err != nil {
		t.Fatal(err)
	}
	rt := ratelimiter.NewDefaultRateLimiter(settings, conn)

	t.Run("Execute by IP", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:8080"

		for i := 0; i < settings.Ratelimit; i++ {
			if proceed, err := rt.Execute(ctx, req); !proceed || err != nil {
				t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, true, nil)
			}
		}

		if proceed, err := rt.Execute(ctx, req); proceed || err != nil {
			t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
		}
	})

	t.Run("Execute by Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
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

	t.Run("Exceed_IP_Rate_Limit", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:8080"
		for i := 0; i < 10; i++ {
			proceed, err := rt.ExecuteByIP(ctx, req)
			if proceed || err != nil {
				t.Errorf("ExecuteByIP(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
			}
		}
	})
}

func TestTokenOverrideIPRateLimit(t *testing.T) {

	ctx := context.Background()
	ipLimit := 10
	tokenLimit := 100
	settings := ratelimiter.NewSettings(ipLimit, 60, true)

	conn, err := db.NewDatabaseConnection(dbOptions)
	if err != nil {
		t.Fatal(err)
	}

	token := "token_with_override"
	if err := conn.CreateToken(ctx, token, tokenLimit); err != nil {
		t.Fatal(err)
	}

	rt := ratelimiter.NewDefaultRateLimiter(settings, conn)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("API_KEY", token)

	for i := 0; i < tokenLimit; i++ {
		if proceed, err := rt.Execute(ctx, req); !proceed || err != nil {
			t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, true, nil)
		}
	}

	if proceed, err := rt.Execute(ctx, req); proceed || err != nil {
		t.Errorf("Execute(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
	}
}

func TestExecuteByIP(t *testing.T) {
	ctx := context.Background()
	token := "g8dXuf2MqNkqJ5tb47qw4m6thqYbrsK24SFZV4OiS83Lmbp8NCYulXtO3tyHJyZN"
	limit := 100
	settings := ratelimiter.NewSettings(10, 60, true)
	conn, err := db.NewDatabaseConnection(dbOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.CreateToken(ctx, token, limit); err != nil {
		t.Fatal(err)
	}
	rt := ratelimiter.NewDefaultRateLimiter(settings, conn)

	t.Run("Execute by IP", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:8080"

		for i := 0; i < settings.Ratelimit; i++ {
			proceed, err := rt.ExecuteByIP(ctx, req)
			if !proceed || err != nil {
				t.Errorf("ExecuteByIP(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, true, nil)
			}
		}

		if proceed, err := rt.ExecuteByIP(ctx, req); proceed || err != nil {
			t.Errorf("ExecuteByIP(%v, %v) = (%v, %v), want (%v, %v)", ctx, req, proceed, err, false, nil)
		}
	})
}
