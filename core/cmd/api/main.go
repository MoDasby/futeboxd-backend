package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-backend/core/internal/user/repository"

	"github.com/modasby/futeboxd-backend/core/database"
	commentHandlers "github.com/modasby/futeboxd-backend/core/internal/comment/handler"
	commentRepository "github.com/modasby/futeboxd-backend/core/internal/comment/repository"
	commentUC "github.com/modasby/futeboxd-backend/core/internal/comment/usecase"

	matchHandlers "github.com/modasby/futeboxd-backend/core/internal/match/handler"
	matchRepository "github.com/modasby/futeboxd-backend/core/internal/match/repository"
	matchUC "github.com/modasby/futeboxd-backend/core/internal/match/usecase"

	profileHandlers "github.com/modasby/futeboxd-backend/core/internal/profile/handler"
	profileRepository "github.com/modasby/futeboxd-backend/core/internal/profile/repository"
	profileUC "github.com/modasby/futeboxd-backend/core/internal/profile/usecase"

	reviewHandlers "github.com/modasby/futeboxd-backend/core/internal/review/handler"
	reviewRepository "github.com/modasby/futeboxd-backend/core/internal/review/repository"
	reviewUC "github.com/modasby/futeboxd-backend/core/internal/review/usecase"

	sessionHandlers "github.com/modasby/futeboxd-backend/core/internal/session/handler"
	sessionRepository "github.com/modasby/futeboxd-backend/core/internal/session/repository"
	sessionUC "github.com/modasby/futeboxd-backend/core/internal/session/usecase"

	userHandlers "github.com/modasby/futeboxd-backend/core/internal/user/handler"
	userUC "github.com/modasby/futeboxd-backend/core/internal/user/usecase"

	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sessionRepo := sessionRepository.NewSessionRepository(db)
	userRepo := repository.NewUserRepository(db)
	profileRepo := profileRepository.NewProfileRepository(db)
	reviewRepo := reviewRepository.NewReviewRepository(db)
	commentRepo := commentRepository.NewCommentsRepository(db)
	matchRepo := matchRepository.NewMatchRepository(db)

	footballClient := football.NewClient()

	sessionUsecases := sessionUC.NewSessionUsecases(sessionRepo, userRepo, footballClient)
	profileUsecases := profileUC.NewProfileUsecases(profileRepo, footballClient)
	reviewUsecases := reviewUC.NewReviewUsecases(reviewRepo, footballClient, userRepo)
	commentUsecases := commentUC.NewCommentUsecases(commentRepo, reviewRepo, userRepo)
	matchUsecases := matchUC.NewMatchUsecases(matchRepo, footballClient)
	userUsecases := userUC.NewUsersUsecases(userRepo, footballClient)

	userHandler := userHandlers.NewUserHandler(userUsecases)
	sessionHandler := sessionHandlers.NewsessionHandler(sessionUsecases)
	profileHandler := profileHandlers.NewProfileHandler(profileUsecases)
	reviewHandler := reviewHandlers.NewReviewsHandler(reviewUsecases)
	commentHandler := commentHandlers.NewCommentHandler(commentUsecases)
	matchHandler := matchHandlers.NewMatchHandler(matchUsecases)

	injectUser := middleware.NewAuthMiddleware(sessionRepo)

	router := http.NewServeMux()

	userHandler.RegisterRoutes(router, injectUser)
	sessionHandler.RegisterRoutes(router, injectUser)
	reviewHandler.RegisterRoutes(router, injectUser)
	profileHandler.RegisterRoutes(router, injectUser)
	commentHandler.RegisterRoutes(router, injectUser)
	matchHandler.RegisterRoutes(router)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
