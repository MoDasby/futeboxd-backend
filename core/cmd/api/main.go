package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/database/aws"
	"github.com/modasby/futeboxd-backend/core/database/postgres"
	"github.com/modasby/futeboxd-backend/core/internal/user/repository"

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

	"github.com/modasby/futeboxd-backend/core/pkg/email"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/log"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panico na rave", "originalError", r)
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := log.InitLogger(cfg.Axiom); err != nil {
		panic(err)
	}

	db, err := postgres.InitDatabase(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s3Client, err := aws.NewS3Client(&cfg.AWS)
	if err != nil {
		panic(err)
	}

	sessionRepo := sessionRepository.NewSessionRepository(db)
	userRepo := repository.NewUserRepository(db)
	profileRepo := profileRepository.NewProfileRepository(db)
	reviewRepo := reviewRepository.NewReviewRepository(db)
	commentRepo := commentRepository.NewCommentsRepository(db)
	matchRepo := matchRepository.NewMatchRepository(db)
	uploadRepo := uploadRepository.NewUploadRepository(s3Client, &cfg.AWS)

	footballClient := football.NewClient(cfg.Football)
	emailClient := email.NewResendClient(cfg.Resend)

	sessionUsecases := sessionUC.NewSessionUsecases(sessionRepo, userRepo)
	profileUsecases := profileUC.NewProfileUsecases(profileRepo, footballClient)
	reviewUsecases := reviewUC.NewReviewUsecases(reviewRepo, footballClient, userRepo)
	commentUsecases := commentUC.NewCommentUsecases(commentRepo, reviewRepo, userRepo)
	matchUsecases := matchUC.NewMatchUsecases(matchRepo, footballClient)
	userUsecases := userUC.NewUsersUsecases(userRepo, footballClient, emailClient, cfg.Frontend)
	uploadUsecases := uploadUC.NewUploadUsecase(uploadRepo, userRepo)

	userHandler := userHandlers.NewUserHandler(userUsecases)
	sessionHandler := sessionHandlers.NewsessionHandler(sessionUsecases, &cfg.Cookies)
	profileHandler := profileHandlers.NewProfileHandler(profileUsecases)
	reviewHandler := reviewHandlers.NewReviewsHandler(reviewUsecases)
	commentHandler := commentHandlers.NewCommentHandler(commentUsecases)
	matchHandler := matchHandlers.NewMatchHandler(matchUsecases)
	uploadHandler := uploadHandlers.NewUploadHandler(uploadUsecases)

	injectUser := middleware.NewAuthMiddleware(sessionRepo, &cfg.Cookies)

	router := http.NewServeMux()

	userHandler.RegisterRoutes(router, injectUser)
	sessionHandler.RegisterRoutes(router, injectUser)
	reviewHandler.RegisterRoutes(router, injectUser)
	profileHandler.RegisterRoutes(router, injectUser)
	commentHandler.RegisterRoutes(router, injectUser)
	matchHandler.RegisterRoutes(router)
	uploadHandler.RegisterRoutes(router, injectUser)

	slog.Debug(fmt.Sprintf("Iniciando servidor na porta: %d", cfg.Server.Port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), middleware.LoggerMiddleware(router.ServeHTTP)); err != nil {
		slog.Error(err.Error())
	}
}
