FROM golang:1.22.5 AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY services/football .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/football/main.go
FROM alpine:latest AS runner

WORKDIR /root/

RUN apk update && apk add postgresql-client

COPY --from=build /app/api ./api
COPY services/football/.env.example .env
COPY services/football/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

ENTRYPOINT [ "sh", "/root/entrypoint.sh" ]