dev:
	@export $(shell grep -v '^#' .env | xargs) && \
	go mod tidy && \
	sh ./wait-to-db-ready.sh && \
	echo "Iniciando a API do football..." && \
	~/go/bin/air --build.cmd "go build -o bin/api cmd/football/main.go" --build.bin "./bin/api"

test:
	@go test ./... -coverprofile=coverage.out

build:
	@go build -o bin/api cmd/football/main.go
