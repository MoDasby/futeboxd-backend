package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/services/users/internal/handler"
	"github.com/modasby/futeboxd-api/services/users/internal/repository"
	"github.com/modasby/futeboxd-api/services/users/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=5434",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %s", err)
	}

	footballClient := football.NewClient("http://localhost:80/football")

	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)

	createUserUseCase := usecase.NewCreateUserUseCase(userRepository, footballClient)
	loginUseCase := usecase.NewLoginUseCase(sessionRepository, userRepository)
	findUserByUsernameUseCase := usecase.NewFindUserByIdOrUsernameUseCase(userRepository, footballClient)
	authenticateUserUseCase := usecase.NewAuthenticateUserUsecase(sessionRepository, userRepository)
	findUserBatchUseCase := usecase.NewFindUserBatchUseCase(userRepository)

	userHandler := handler.NewUserHandler(createUserUseCase, loginUseCase, findUserByUsernameUseCase, authenticateUserUseCase, findUserBatchUseCase)

	r := http.NewServeMux()

	r.HandleFunc("POST /auth/register", userHandler.CreateUser)
	r.HandleFunc("POST /auth/login", userHandler.Login)
	r.HandleFunc("GET /auth/user", userHandler.AuthenticateUser)
	r.HandleFunc("/user/{identificator}", userHandler.FindByIdOrUsername)
	r.HandleFunc("POST /batch/user", userHandler.FindBatchByID)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
