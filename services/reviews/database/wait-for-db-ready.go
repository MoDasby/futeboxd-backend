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

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", os.Getenv("db_user"), os.Getenv("db_name"), os.Getenv("db_password"), os.Getenv("db_host"))
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
