package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDatabase() (*sql.DB, error) {
	envPath, err := filepath.Abs(".env")
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(envPath); err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db, nil
}
