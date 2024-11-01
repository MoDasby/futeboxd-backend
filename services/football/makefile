dev:
	@go mod tidy
	@sh ../../scripts/wait-to-db-ready.sh --env-file ../../env/football.env
	@echo "Iniciando a API do football..."
	@~/go/bin/air --build.cmd "go build -o bin/api cmd/football/main.go" --build.bin "./bin/api"

test:
	@go test ./... -coverprofile=coverage.out

build:
	@go build -o bin/api cmd/football/main.go

compose:down:
	@docker compose --env-file ../../env/football.env -f ../infra/football/docker/compose.dev.yml down

compose:up:
	@docker compose --env-file ../../env/football.env -f ../infra/football/docker/compose.dev.yml up -d

compose:stop:
	@docker compose --env-file ../../env/football.env -f ../infra/football/docker/compose.dev.yml stop
