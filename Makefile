build:
	docker-compose --env-file dev.env up -d --build
	@echo "Docker compose executado com sucesso."

up:
	docker-compose --env-file dev.env up -d
	@echo "Docker compose executado com sucesso."

down:
	docker-compose down
	@echo "Contêineres do aplicativo parados e removidos."

logs:
	docker logs rate-limiter
	@echo "Logs carregados com sucesso."

test-run:
	go test -v -run TestExecutionTheRateLimiter ./internal/ratelimiter/default_test.go
	@echo "Testes realizados com sucesso."

test-exceed:
	go test -v -run TestExceedingTheRateLimiter ./internal/ratelimiter/default_test.go
	@echo "Testes realizados com sucesso."

test-override:
	go test -v -run TestTokenOverrideIPRateLimit ./internal/ratelimiter/default_test.go
	@echo "Testes realizados com sucesso."

test-ip:
	go test -v -run TestExecuteByIP ./internal/ratelimiter/default_test.go
	@echo "Testes realizados com sucesso."

.PHONY: go	