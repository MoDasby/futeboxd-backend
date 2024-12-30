package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	uploadHandlers "github.com/modasby/futeboxd-backend/core/internal/upload"
	uploadRepository "github.com/modasby/futeboxd-backend/core/internal/upload/repository"
	uploadUC "github.com/modasby/futeboxd-backend/core/internal/upload/usecase"

	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy-sla", "dummy-secret", ""),
		),
	)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.UsePathStyle = true
	})

	sessionRepo := sessionRepository.NewSessionRepository(db)
	userRepo := repository.NewUserRepository(db)
	profileRepo := profileRepository.NewProfileRepository(db)
	reviewRepo := reviewRepository.NewReviewRepository(db)
	commentRepo := commentRepository.NewCommentsRepository(db)
	matchRepo := matchRepository.NewMatchRepository(db)
	uploadRepo := uploadRepository.NewUploadRepository(client)

	footballClient := football.NewClient()

	sessionUsecases := sessionUC.NewSessionUsecases(sessionRepo, userRepo, footballClient)
	profileUsecases := profileUC.NewProfileUsecases(profileRepo, footballClient)
	reviewUsecases := reviewUC.NewReviewUsecases(reviewRepo, footballClient, userRepo)
	commentUsecases := commentUC.NewCommentUsecases(commentRepo, reviewRepo, userRepo)
	matchUsecases := matchUC.NewMatchUsecases(matchRepo, footballClient)
	userUsecases := userUC.NewUsersUsecases(userRepo, footballClient)
	uploadUsecases := uploadUC.NewUploadUsecase(uploadRepo, userRepo)

	userHandler := userHandlers.NewUserHandler(userUsecases)
	sessionHandler := sessionHandlers.NewsessionHandler(sessionUsecases)
	profileHandler := profileHandlers.NewProfileHandler(profileUsecases)
	reviewHandler := reviewHandlers.NewReviewsHandler(reviewUsecases)
	commentHandler := commentHandlers.NewCommentHandler(commentUsecases)
	matchHandler := matchHandlers.NewMatchHandler(matchUsecases)
	uploadHandler := uploadHandlers.NewUploadHandler(uploadUsecases)

	injectUser := middleware.NewAuthMiddleware(sessionRepo)

	router := http.NewServeMux()

	userHandler.RegisterRoutes(router, injectUser)
	sessionHandler.RegisterRoutes(router, injectUser)
	reviewHandler.RegisterRoutes(router, injectUser)
	profileHandler.RegisterRoutes(router, injectUser)
	commentHandler.RegisterRoutes(router, injectUser)
	matchHandler.RegisterRoutes(router)
	uploadHandler.RegisterRoutes(router, injectUser)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
