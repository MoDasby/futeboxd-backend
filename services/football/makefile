dev:
	@go mod tidy

	@sh ./wait-to-db-ready.sh

	@echo "Iniciando a API em Go..."
	@~/go/bin/air --build.cmd "go build -o bin/api cmd/football/main.go" --build.bin "./bin/api"

test:
	@go test ./... -coverprofile=coverage.out

build:
	@go build -o bin/api cmd/football/main.go

compose\:down:
	@docker compose --env-file .env -f database/compose.dev.yml down

compose\:up:
	@docker compose --env-file .env -f database/compose.dev.yml up -d

compose\:stop:
	@docker compose --env-file .env -f database/compose.dev.yml stop
