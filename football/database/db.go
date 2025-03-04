package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func InitDatabase() (*sql.DB, error) {
	DBPASS := os.Getenv("POSTGRES_PASSWORD")
	DBUSER := os.Getenv("DB_USER")
	DBNAME := os.Getenv("POSTGRES_DB")
	DBHOST := os.Getenv("DB_HOST")
	PGPORT := os.Getenv("PGPORT")

	if os.Getenv("ENV") == "prod" {
		secretPassword, err := os.ReadFile(os.Getenv("DB_PASSWORD_FILE"))
		if err != nil {
			return nil, err
		}

		DBPASS = strings.TrimSpace(string(secretPassword))
	}

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=%s",
		DBUSER,
		DBNAME,
		DBPASS,
		DBHOST,
		PGPORT,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
