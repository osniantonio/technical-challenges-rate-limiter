up:
	docker-compose --env-file dev.env up -d --build
	@echo "Docker compose executado com sucesso."

down:
	docker-compose down
	@echo "Contêineres do aplicativo parados e removidos."

test:
	go test -v ./...
	@echo "Testes realizados com sucesso."

coverage:
	go test -v ./... -coverprofile=c.out
	go tool cover -html=c.out
	@echo "Cobertura de testes realizada com sucesso."

.PHONY: go	