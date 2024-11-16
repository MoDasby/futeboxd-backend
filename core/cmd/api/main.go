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
	profileRepo := repository.NewProfileRepository(db)

	footballClient := football.NewClient()

	listTrendingMatchesQuery := queries.NewListTrendingMatchQuery(db, footballClient)

	createReviewUseCase := reviewsUsecases.NewCreateReviewUseCase(reviewRepository, footballClient)
	listReviewsUseCase := reviewsUsecases.NewListReviewsUseCase(reviewRepository, footballClient)
	listFeedUsecase := reviewsUsecases.NewListFeedUsecase(reviewRepository, userRepository, footballClient)
	deleteReviewsUseCase := reviewsUsecases.NewDeleteReviewUseCase(reviewRepository)
	createCommentUsecase := reviewsUsecases.NewCreateCommentUsecase(commentsRepo, reviewRepository)
	listCommentsUsecase := reviewsUsecases.NewGetCommentsUsecase(commentsRepo)
	toggleLikeCommentUsecase := reviewsUsecases.NewToggleLikeCommentUsecase(commentsRepo)
	toggleLikeReviewUsecase := reviewsUsecases.NewToggleLikeReviewUsecase(reviewRepository)
	deleteCommentUsecase := reviewsUsecases.NewDeleteCommentUsecase(commentsRepo)

	createUserUseCase := usersUsecases.NewCreateUserUseCase(userRepository, footballClient)
	findByProfile := usersUsecases.NewFindProfileUsecase(profileRepo, footballClient)
	editUserUsecase := usersUsecases.NewEditUserUseCase(userRepository, footballClient)
	loginUsecase := usersUsecases.NewLoginUseCase(sessionRepository, userRepository)
	logoutUsecase := usersUsecases.NewLogoutUsecase(sessionRepository)
	updatePasswordUsecase := usersUsecases.NewUpdatePasswordUsecase(userRepository)
	toggleFollowUsecase := usersUsecases.NewToggleFollowUsecase(userRepository, followersRepo)
	searchProfile := usersUsecases.NewSearchProfile(profileRepo, footballClient)

	reviewHandler := handler.NewReviewHandler(
		createReviewUseCase,
		listReviewsUseCase,
		listFeedUsecase,
		deleteReviewsUseCase,
		createCommentUsecase,
		listCommentsUsecase,
		toggleLikeCommentUsecase,
		toggleLikeReviewUsecase,
		deleteCommentUsecase,
		listTrendingMatchesQuery,
	)
	userHandler := handler.NewUserHandler(
		createUserUseCase,
		findByProfile,
		editUserUsecase,
		loginUsecase,
		logoutUsecase,
		updatePasswordUsecase,
		toggleFollowUsecase,
		searchProfile,
	)

	injectUser := middleware.NewInjectUserMiddleware(sessionRepository, userRepository)

	router := http.NewServeMux()

	reviewHandler.RegisterRoutes(router, injectUser)
	userHandler.RegisterRoutes(router, injectUser)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
