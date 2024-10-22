package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/pkg/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
	usecase "github.com/modasby/futeboxd-api/services/core/internal/usecase/reviews"
)

type ReviewHandler struct {
	createReviewUseCase *usecase.CreateReviewUseCase
	listReviewsUseCase  *usecase.ListReviewsUseCase
	listFeedUsecase     *usecase.ListFeedUsecase
	deleteReviewUseCase *usecase.DeleteReviewUseCase
}

func NewReviewHandler(
	createReviewUseCase *usecase.CreateReviewUseCase,
	listReviewsUseCase *usecase.ListReviewsUseCase,
	listFeedUsecase *usecase.ListFeedUsecase,
	deleteReviewUseCase *usecase.DeleteReviewUseCase,
) *ReviewHandler {
	return &ReviewHandler{
		createReviewUseCase: createReviewUseCase,
		listReviewsUseCase:  listReviewsUseCase,
		listFeedUsecase:     listFeedUsecase,
		deleteReviewUseCase: deleteReviewUseCase,
	}
}

func (h *ReviewHandler) RegisterRoutes(router *http.ServeMux, injectUser middleware.Middleware) {
	router.HandleFunc("POST /reviews", injectUser(h.create))
	router.HandleFunc("/reviews/feed", injectUser(h.listFeed))
	router.HandleFunc("/reviews", h.listAll)
	router.HandleFunc("DELETE /reviews/delete/{reviewID}", injectUser(h.delete))
}

func (h *ReviewHandler) create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil || user.Username == "anonymous" {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	var input usecase.CreateReviewInputDTO

	input.Author = user

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpErr := errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:REVIEW:CREATE:INVALID_BODY",
		)

		errors.HandleHttpError(w, httpErr)

		return
	}

	if err := h.createReviewUseCase.Execute(input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *ReviewHandler) listAll(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	team := r.URL.Query().Get("team")
	match := r.URL.Query().Get("match")

	size, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		size = 50
	}

	index, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		index = 1
	}

	page, err := pagination.WithPage(int(size), int(index))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	reviews, err := h.listReviewsUseCase.Execute(username, team, match, *page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil || user.Username == "anonymous" {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	reviewIDStr := r.PathValue("reviewID")
	reviewID, err := strconv.ParseInt(reviewIDStr, 10, 64)
	if err != nil {
		err = errors.NewHTTPErr(
			"review id está em um formato inválido",
			400,
			"HANDLER:REVIEW:DELETE:REVIEW_ID_INVALID_FORMAT",
		)

		errors.HandleHttpError(w, err)

		return
	}

	if err := h.deleteReviewUseCase.Execute(reviewID, user.ID); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) listFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil || user.Username == "anonymous" {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	strategy := r.URL.Query().Get("strategy")

	size, err := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		size = 50
	}

	index, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		index = 1
	}

	page, err := pagination.WithPage(int(size), int(index))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	reviews, err := h.listFeedUsecase.Execute(user.ID, strategy, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
