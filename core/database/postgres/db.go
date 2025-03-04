package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/modasby/futeboxd-backend/core/config"
)

func InitDatabase(cfg config.Postgres) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=%d",
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	db, err := waitToDbReady(connStr, 5, time.Second*3)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func waitToDbReady(connStr string, maxAttempts int, delay time.Duration) (*sql.DB, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			slog.Error(
				"Erro ao abrir conexão com o banco de dados (tentativa %d/%d): %v\n",
				"attempt", attempt,
				"maxAttempts", maxAttempts,
				"originalError", err.Error(),
			)
			time.Sleep(delay)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			log.Printf("Erro ao conectar ao banco de dados (tentativa %d/%d): %v\n", attempt, maxAttempts, err)
			db.Close()
			time.Sleep(delay)
			continue
		}

		return db, nil
	}

	return nil, errors.New("não foi possível conectar com o banco de dados")
}
