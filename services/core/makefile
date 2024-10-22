test:
	@go test ./... -coverprofile=coverage.out

dev:
	@go mod tidy
	
	@echo "Aguardando o banco de dados ficar disponível..."
	@./wait-to-db-ready.sh

	@echo "Iniciando a API em Go..."
	@~/go/bin/air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

build:
	@go build -o bin/api cmd/api/main.go

compose\:down:
	@docker compose --env-file .env -f database/compose.dev.yml down

compose\:up:
	@docker compose --env-file .env -f database/compose.dev.yml up -d

compose\:stop:
	@docker compose --env-file .env -f database/compose.dev.yml stop
