package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/services/core/database"
	"github.com/modasby/futeboxd-api/services/core/internal/handler"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
	"github.com/modasby/futeboxd-api/services/core/internal/queries"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	reviewsUsecases "github.com/modasby/futeboxd-api/services/core/internal/usecase/reviews"
	usersUsecases "github.com/modasby/futeboxd-api/services/core/internal/usecase/users"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	reviewRepository := repository.NewReviewRepository(db)
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	followersRepo := repository.NewFollowersRepository(db)
	commentsRepo := repository.NewCommentsRepository(db)

	footballClient := football.NewClient("http://localhost:80/football")

	followStatsQueryService := queries.NewFollowStatsQuery(db)

	createReviewUseCase := reviewsUsecases.NewCreateReviewUseCase(reviewRepository, footballClient)
	listReviewsUseCase := reviewsUsecases.NewListReviewsUseCase(reviewRepository, userRepository)
	listFeedUsecase := reviewsUsecases.NewListFeedUsecase(reviewRepository, userRepository)
	deleteReviewsUseCase := reviewsUsecases.NewDeleteReviewUseCase(reviewRepository)
	createCommentUsecase := reviewsUsecases.NewCreateCommentUsecase(commentsRepo, userRepository)
	listCommentsUsecase := reviewsUsecases.NewGetCommentsUsecase(commentsRepo)

	createUserUseCase := usersUsecases.NewCreateUserUseCase(userRepository, footballClient)
	findByProfile := usersUsecases.NewFindProfile(userRepository, footballClient, followStatsQueryService)
	editUserUsecase := usersUsecases.NewEditUserUseCase(userRepository)
	loginUsecase := usersUsecases.NewLoginUseCase(sessionRepository, userRepository)
	followUserUsecase := usersUsecases.NewFollowUserUseCase(followersRepo, userRepository)
	unfollowUserUsecase := usersUsecases.NewUnfollowUserUseCase(followersRepo, userRepository)

	reviewHandler := handler.NewReviewHandler(
		createReviewUseCase,
		listReviewsUseCase,
		listFeedUsecase,
		deleteReviewsUseCase,
		createCommentUsecase,
		listCommentsUsecase,
	)
	userHandler := handler.NewUserHandler(
		createUserUseCase,
		findByProfile,
		editUserUsecase,
		loginUsecase,
		followUserUsecase,
		unfollowUserUsecase,
	)

	injectUser := middleware.NewInjectUserMiddleware(sessionRepository, userRepository)

	router := http.NewServeMux()

	reviewHandler.RegisterRoutes(router, injectUser)
	userHandler.RegisterRoutes(router, injectUser)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
