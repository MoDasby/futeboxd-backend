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
	"github.com/modasby/futeboxd-api/pkg/client/user"
	"github.com/modasby/futeboxd-api/services/reviews/internal/handler"
	"github.com/modasby/futeboxd-api/services/reviews/internal/middleware"
	"github.com/modasby/futeboxd-api/services/reviews/internal/repository"
	"github.com/modasby/futeboxd-api/services/reviews/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", os.Getenv("db_user"), os.Getenv("db_name"), os.Getenv("db_password"), os.Getenv("db_host"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %s", err)
	}

	footballClient := football.NewClient("http://localhost:80/football")
	userClient := user.NewClient("http://localhost:8082")

	reviewRepository := repository.NewReviewRepository(db)

	createReviewUseCase := usecase.NewCreateReviewUseCase(reviewRepository, footballClient, userClient)
	listReviewsUseCase := usecase.NewListReviewsUseCase(reviewRepository, footballClient, userClient)

	reviewHandler := handler.NewReviewHandler(createReviewUseCase, listReviewsUseCase)

	injectUserMiddleware := middleware.NewInjectUserMiddleware(userClient)

	r := http.NewServeMux()

	r.HandleFunc("POST /review", injectUserMiddleware(reviewHandler.CreateReview))
	r.HandleFunc("/review", reviewHandler.ListReviews)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
