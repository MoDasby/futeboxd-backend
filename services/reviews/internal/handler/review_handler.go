package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/pkg/client/user"
	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/reviews/internal/middleware"
	"github.com/modasby/futeboxd-api/services/reviews/internal/usecase"
)

type ReviewHandler struct {
	createReviewUseCase *usecase.CreateReviewUseCase
	listReviewsUseCase  *usecase.ListReviewsUseCase
}

func NewReviewHandler(
	createReviewUseCase *usecase.CreateReviewUseCase,
	listReviewsUseCase *usecase.ListReviewsUseCase,
) *ReviewHandler {
	return &ReviewHandler{
		createReviewUseCase: createReviewUseCase,
		listReviewsUseCase:  listReviewsUseCase,
	}
}

func (handler *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(middleware.SessionKey).(*user.AuthenticateOutputDTO)
	if !ok || session == nil {
		httpErr := errors.NewErrUnauthorized("ocorreu um erro de autenticação, tente logar novamente")
		errors.HandleHttpError(w, httpErr)

		return
	}

	var review usecase.CreateReviewInputDTO
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		httpErr := errors.NewErrBadRequest("Corpo de requisição inválido")
		errors.HandleHttpError(w, httpErr)

		return
	}

	review.UserID = session.User.ID
	review.AuthorUsername = session.User.Username

	err = handler.createReviewUseCase.Execute(review)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (handler *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	size := 50
	user := r.URL.Query().Get("user")
	team := r.URL.Query().Get("team")
	match := r.URL.Query().Get("match")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	output, err := handler.listReviewsUseCase.Execute(size, pageInt, user, team, match)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
