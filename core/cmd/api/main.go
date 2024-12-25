package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-backend/core/internal/user/repository"

	"github.com/modasby/futeboxd-backend/core/database"
	"github.com/modasby/futeboxd-backend/core/internal/comment"
	commentRepository "github.com/modasby/futeboxd-backend/core/internal/comment/repository"
	"github.com/modasby/futeboxd-backend/core/internal/match"
	matchRepository "github.com/modasby/futeboxd-backend/core/internal/match/repository"
	"github.com/modasby/futeboxd-backend/core/internal/profile"
	profileRepository "github.com/modasby/futeboxd-backend/core/internal/profile/repository"
	reviews "github.com/modasby/futeboxd-backend/core/internal/review"
	reviewRepository "github.com/modasby/futeboxd-backend/core/internal/review/repository"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	sessionRepository "github.com/modasby/futeboxd-backend/core/internal/session/repository"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/auth"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
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

	sessionUsecases := session.NewSessionUsecases(sessionRepo, userRepo, footballClient)
	profileUsecases := profile.NewProfileUsecases(profileRepo, footballClient)
	reviewUsecases := reviews.NewReviewUsecases(reviewRepo, footballClient, userRepo)
	commentUsecases := comment.NewCommentUsecases(commentRepo, reviewRepo, userRepo)
	matchUsecases := match.NewMatchUsecases(matchRepo, footballClient)
	userUsecases := user.NewUsersUsecases(userRepo, footballClient)

	userHandler := user.NewUserHandler(userUsecases)
	sessionHandler := session.NewsessionHandler(sessionUsecases)
	profileHandler := profile.NewProfileHandler(profileUsecases)
	reviewHandler := reviews.NewReviewsHandler(reviewUsecases)
	commentHandler := comment.NewCommentHandler(commentUsecases)
	matchHandler := match.NewMatchHandler(matchUsecases)

	injectUser := auth.NewAuthMiddleware(sessionRepo)

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
