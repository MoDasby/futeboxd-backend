package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/core/database"
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/handler"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
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
	profileRepo := repository.NewProfileRepository(db)

	footballClient := football.NewClient()

	createReviewUseCase := reviewsUsecases.NewCreateReviewUseCase(reviewRepository, footballClient)
	listReviewsUseCase := reviewsUsecases.NewListReviewsUseCase(reviewRepository, userRepository)
	listFeedUsecase := reviewsUsecases.NewListFeedUsecase(reviewRepository, userRepository)
	deleteReviewsUseCase := reviewsUsecases.NewDeleteReviewUseCase(reviewRepository)
	createCommentUsecase := reviewsUsecases.NewCreateCommentUsecase(commentsRepo, reviewRepository)
	listCommentsUsecase := reviewsUsecases.NewGetCommentsUsecase(commentsRepo)
	toggleLikeCommentUsecase := reviewsUsecases.NewToggleLikeCommentUsecase(commentsRepo)
	toggleLikeReviewUsecase := reviewsUsecases.NewToggleLikeReviewUsecase(reviewRepository)

	createUserUseCase := usersUsecases.NewCreateUserUseCase(userRepository, footballClient)
	findByProfile := usersUsecases.NewFindProfileUsecase(profileRepo, footballClient)
	editUserUsecase := usersUsecases.NewEditUserUseCase(userRepository, footballClient)
	loginUsecase := usersUsecases.NewLoginUseCase(sessionRepository, userRepository)
	logoutUsecase := usersUsecases.NewLogoutUsecase(sessionRepository)
	updatePasswordUsecase := usersUsecases.NewUpdatePasswordUsecase(userRepository)
	toggleFollowUsecase := usersUsecases.NewToggleFollowUsecase(userRepository, followersRepo)

	reviewHandler := handler.NewReviewHandler(
		createReviewUseCase,
		listReviewsUseCase,
		listFeedUsecase,
		deleteReviewsUseCase,
		createCommentUsecase,
		listCommentsUsecase,
		toggleLikeCommentUsecase,
		toggleLikeReviewUsecase,
	)
	userHandler := handler.NewUserHandler(
		createUserUseCase,
		findByProfile,
		editUserUsecase,
		loginUsecase,
		logoutUsecase,
		updatePasswordUsecase,
		toggleFollowUsecase,
	)

	injectUser := middleware.NewInjectUserMiddleware(sessionRepository, userRepository)

	router := http.NewServeMux()

	reviewHandler.RegisterRoutes(router, injectUser)
	userHandler.RegisterRoutes(router, injectUser)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
