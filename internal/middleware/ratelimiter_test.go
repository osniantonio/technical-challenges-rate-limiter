package middleware_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/handler"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/infra/db"
	middlewareTest "github.com/osniantonio/technical-challenges-rate-limiter/internal/middleware"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/ratelimiter"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbOptions                  *gateway.DatabaseOptions
	tooManyRequestsMessageTest = "you have reached the maximum number of requests or actions allowed within a certain time frame"
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

func TestRateLimiterMiddleware(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to Docker: %v", err)
	}

	resource, err := pool.Run("redis", "latest", nil)
	if err != nil {
		t.Fatalf("Could not start Redis container: %v", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge Redis container: %v", err)
		}
	}()

	// Redis
	dbOptions := &gateway.DatabaseOptions{
		Protocol: "redis",
		Host:     resource.Container.NetworkSettings.IPAddress,
		Port:     "6379",
		Database: "0",
		Password: "",
	}

	// Conecte-se ao banco de dados Redis
	conn, err := db.NewDatabaseConnection(dbOptions)
	if err != nil {
		t.Fatalf("Could not connect to Redis: %v", err)
	}

	// Configurar o rate limiter
	settings := ratelimiter.NewSettings(10, 60, true)
	rt := ratelimiter.NewDefaultRateLimiter(settings, conn)

	mid := middlewareTest.NewRateLimiterMiddleware(rt)

	ctx := context.Background()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.RemoteAddr = "127.0.0.1:8080"

	// Executar o teste
	t.Run("Execute", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			h := mid.Execute(ctx, &handler.DefaultHandler{})
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			expected := "Welcome"
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
			}
		}
		h := mid.Execute(ctx, &handler.DefaultHandler{})
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusTooManyRequests {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
		}
		if rr.Body.String() != tooManyRequestsMessageTest {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tooManyRequestsMessageTest)
		}
	})
}
