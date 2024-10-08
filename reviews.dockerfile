FROM golang:1.22.5 AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY services/reviews .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api/main.go
FROM alpine:latest AS runner

WORKDIR /root/

RUN apk update && apk add postgresql-client

COPY --from=build /app/api ./api
COPY services/reviews/.env.example .env
COPY services/reviews/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

ENTRYPOINT [ "sh", "/root/entrypoint.sh" ]