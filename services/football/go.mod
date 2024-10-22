module github.com/modasby/futeboxd-api/services/football

go 1.22.5

replace github.com/modasby/futeboxd-api/pkg => ../../pkg

require (
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/modasby/futeboxd-api/pkg v0.0.0-00010101000000-000000000000
)
