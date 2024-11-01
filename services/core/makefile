dev:
	@go mod tidy
	@sh ../../scripts/wait-to-db-ready.sh --env-file ../../env/core.env
	@echo "Iniciando a API do core..."
	@~/go/bin/air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

test:
	@go test ./... -coverprofile=coverage.out

build:
	@go build -o bin/api cmd/api/main.go

compose:down:
	@docker compose --env-file ../../env/core.env -f ../infra/core/docker/compose.dev.yml down

compose:up:
	@docker compose --env-file ../../env/core.env -f ../infra/core/docker/compose.dev.yml up -d

compose:stop:
	@docker compose --env-file ../../env/core.env -f ../infra/core/docker/compose.dev.yml stop
