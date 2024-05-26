package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/osniantonio/technical-challenges-rate-limiter/configs"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/gateway"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/handler"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/infra/db"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/middleware"
	"github.com/osniantonio/technical-challenges-rate-limiter/internal/ratelimiter"
)

type Token struct {
	Token string `json:"token"`
	Limit int    `json:"limit"`
}

var (
	config   *configs.Config
	settings *ratelimiter.Settings
	conn     gateway.DatabaseGateway
	rt       ratelimiter.RateLimiter
	mid      middleware.Middleware
)

func loadConfig() {
	var err error
	config, err = configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
}

func NewDatabaseConnection() {
	var err error
	conn, err = db.NewDatabaseConnection(&gateway.DatabaseOptions{
		Protocol: config.DBProtocol,
		Host:     config.DBHost,
		Port:     config.DBPort,
		Password: config.DBPassword,
		Database: config.DBDatabase,
	})
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func NewTokens() {
	r, err := os.ReadFile("./resources/tokens.json")
	if err != nil {
		log.Fatalf("Failed when trying to load the tokens: %v", err)
	}
	tokens := []Token{}
	err = json.Unmarshal([]byte(r), &tokens)
	if err != nil {
		log.Fatalf("Failed when trying to parse the tokens: %v", err)
	}
	ctx := context.Background()
	for _, token := range tokens {
		conn.CreateToken(ctx, token.Token, token.Limit)
	}
}

func init() {
	loadConfig()
	NewDatabaseConnection()
	NewTokens()
	settings = ratelimiter.NewSettings(config.RateLimit, config.ExpirationTime, config.LimitByToken)
	rt = ratelimiter.NewDefaultRateLimiter(settings, conn)
	mid = middleware.NewRateLimiterMiddleware(rt)
}

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	mux.Handle("/", mid.Execute(ctx, &handler.DefaultHandler{}))
	http.ListenAndServe(":8080", mux)
}
