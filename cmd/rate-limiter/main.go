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

func main() {
	loadConfig()
	connectToDatabase()
	createTokens()
	initSettings()
	initRateLimiter()
	initMiddleware()
	startServer()
}

func loadConfig() {
	var err error
	config, err = configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
}

func connectToDatabase() {
	var err error
	conn, err = db.NewDatabaseConnection(&gateway.DatabaseOptions{
		Protocol: config.DBProtocol,
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		Database: config.DBDatabase,
	})
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func createTokens() {
	r, err := os.ReadFile("./assets/tokens.json")
	if err != nil {
		log.Fatalf("unable to load tokens: %v", err)
	}
	tokens := []Token{}
	err = json.Unmarshal([]byte(r), &tokens)
	if err != nil {
		log.Fatalf("unable to parse tokens: %v", err)
	}
	ctx := context.Background()
	for _, token := range tokens {
		conn.CreateToken(ctx, token.Token, token.Limit)
	}
}

func initSettings() {
	settings = ratelimiter.NewSettings(config.RateLimit, config.ExpirationTime, config.LimitByToken)
}

func initRateLimiter() {
	rt = ratelimiter.NewDefaultRateLimiter(settings, conn)
}

func initMiddleware() {
	mid = middleware.NewRateLimiterMiddleware(rt)
}

func startServer() {
	ctx := context.Background()
	mux := http.NewServeMux()
	mux.Handle("/", mid.Execute(ctx, &handler.DefaultHandler{}))
	http.ListenAndServe(":8080", mux)
}
