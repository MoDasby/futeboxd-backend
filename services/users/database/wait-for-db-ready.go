package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	logger := log.Default()

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	attempts := 0
	ready := false

	for attempts <= 5 && !ready {
		time.Sleep(time.Second * 1)

		err = db.Ping()
		if err != nil {
			log.Print(err)
			logger.Println("Banco de dados não está pronto, tentando novamente...")
		} else {
			ready = true
			break
		}

		attempts++
	}

	if ready {
		logger.Println("Banco de dados está pronto.")

		return
	}

	logger.Fatal("Falha no banco de dados")
}
