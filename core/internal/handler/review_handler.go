package handler

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-api/services/core/internal/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/queries"
	usecase "github.com/modasby/futeboxd-api/services/core/internal/usecase/reviews"
)

type ReviewHandler struct {
	createReviewUseCase      *usecase.CreateReviewUseCase
	listReviewsUseCase       *usecase.ListReviewsUseCase
	listFeedUsecase          *usecase.ListFeedUsecase
	deleteReviewUseCase      *usecase.DeleteReviewUseCase
	createCommentUsecase     *usecase.CreateCommentUsecase
	listCommentsUsecase      *usecase.GetCommentsUsecase
	toggleLikeCommentUsecase *usecase.ToggleLikeCommentUsecase
	toggleLikeReviewUsecase  *usecase.ToggleLikeReviewUsecase
	deleteCommentUsecase     *usecase.DeleteCommentUsecase
	listTrendingMatches      *queries.ListTrendingMatchQuery
}

func NewReviewHandler(
	createReviewUseCase *usecase.CreateReviewUseCase,
	listReviewsUseCase *usecase.ListReviewsUseCase,
	listFeedUsecase *usecase.ListFeedUsecase,
	deleteReviewUseCase *usecase.DeleteReviewUseCase,
	createCommentUsecase *usecase.CreateCommentUsecase,
	listCommentsUsecase *usecase.GetCommentsUsecase,
	toggleLikeCommentUsecase *usecase.ToggleLikeCommentUsecase,
	toggleLikeReviewUsecase *usecase.ToggleLikeReviewUsecase,
	deleteCommentUsecase *usecase.DeleteCommentUsecase,
	listTrendingMatches *queries.ListTrendingMatchQuery,
) *ReviewHandler {
	return &ReviewHandler{
		createReviewUseCase:      createReviewUseCase,
		listReviewsUseCase:       listReviewsUseCase,
		listFeedUsecase:          listFeedUsecase,
		deleteReviewUseCase:      deleteReviewUseCase,
		createCommentUsecase:     createCommentUsecase,
		listCommentsUsecase:      listCommentsUsecase,
		toggleLikeCommentUsecase: toggleLikeCommentUsecase,
		toggleLikeReviewUsecase:  toggleLikeReviewUsecase,
		deleteCommentUsecase:     deleteCommentUsecase,
		listTrendingMatches:      listTrendingMatches,
	}
}

func (h *ReviewHandler) RegisterRoutes(router *http.ServeMux, injectUser middleware.Middleware) {
	router.HandleFunc("POST /reviews", injectUser(h.create, false))
	router.HandleFunc("GET /reviews/feed", injectUser(h.listFeed, false))
	router.HandleFunc("GET /reviews", injectUser(h.listAll, true))
	router.HandleFunc("DELETE /reviews/{reviewID}", injectUser(h.delete, false))
	router.HandleFunc("POST /reviews/{reviewID}/comments", injectUser(h.createComment, false))
	router.HandleFunc("DELETE /reviews/{reviewID}/comments/{commentID}", injectUser(h.deleteComment, false))
	router.HandleFunc("GET /reviews/{reviewID}/comments", injectUser(h.listComments, true))
	router.HandleFunc("POST /reviews/{reviewID}/like", injectUser(h.toggleLikeReview, false))
	router.HandleFunc("POST /reviews/{reviewID}/comments/{commentID}/like", injectUser(h.toggleLikeComment, false))
	router.HandleFunc("GET /matches/trending", injectUser(h.getTrending, true))
}

func (h *ReviewHandler) getTrending(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.listTrendingMatches.Execute(page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := sendJsonResponse(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) deleteComment(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	commentID, err := ParseIntValue(r.PathValue("commentID"))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := h.deleteCommentUsecase.Execute(user.ID, commentID); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) toggleLikeReview(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	reviewID, err := ParseIntValue(r.PathValue("reviewID"))
	if err != nil {
		err = errors.NewHTTPErr(
			"id de reviews só pode conter números positivos",
			400,
			"HANDLER:LIKE_REVIEW:INVALID_REVIEW_ID_FORMAT",
		)

		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.toggleLikeReviewUsecase.Execute(user.ID, reviewID)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := sendJsonResponse(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) toggleLikeComment(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	commentID, err := ParseIntValue(r.PathValue("commentID"))
	if err != nil {
		err = errors.NewHTTPErr(
			"id de comentarios só pode conter números positivos",
			400,
			"HANDLER:LIKE_COMMENT:INVALID_COMMENT_ID_FORMAT",
		)

		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.toggleLikeCommentUsecase.Execute(user.ID, commentID)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := sendJsonResponse(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) listComments(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	reviewID, err := ParseIntValue(r.PathValue("reviewID"))
	if err != nil {
		httpErr := errors.NewHTTPErr(
			"id inválido para a review",
			400,
			"HANDLER:REVIEW:CREATE_COMMENT:INVALID_REVIEW_ID",
		)

		errors.HandleHttpError(w, httpErr)

		return
	}

	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.listCommentsUsecase.Execute(user.ID, reviewID, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := sendJsonResponse(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) createComment(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	reviewID, err := ParseIntValue(r.PathValue("reviewID"))
	if err != nil {
		httpErr := errors.NewHTTPErr(
			"id inválido para a review",
			400,
			"HANDLER:REVIEW:CREATE_COMMENT:INVALID_REVIEW_ID",
		)

		errors.HandleHttpError(w, httpErr)

		return
	}

	var input usecase.CommentInputDTO

	input.ParentID = reviewID

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpErr := errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:REVIEW:CREATE_COMMENT:INVALID_BODY",
		)

		errors.HandleHttpError(w, httpErr)

		return
	}

	input.Author = user

	if err := h.createCommentUsecase.Execute(input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) create(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	var input usecase.CreateReviewInputDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpErr := errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:REVIEW:CREATE:INVALID_BODY",
		)

		errors.HandleHttpError(w, httpErr)

		return
	}

	input.Author = user

	if err := h.createReviewUseCase.Execute(input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) listAll(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	username := r.URL.Query().Get("username")
	team := r.URL.Query().Get("team")
	match := r.URL.Query().Get("match")
	searchTerm := r.URL.Query().Get("search")

	options := usecase.ListReviewsOptions{
		Username:   username,
		Team:       team,
		Match:      match,
		SearchTerm: searchTerm,
	}

	output, err := h.listReviewsUseCase.Execute(user.ID, options, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := sendJsonResponse(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) delete(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	reviewID, err := ParseIntValue(r.PathValue("reviewID"))
	if err != nil {
		err = errors.NewHTTPErr(
			"review id está em um formato inválido",
			400,
			"HANDLER:REVIEW:DELETE:REVIEW_ID_INVALID_FORMAT",
		)

		errors.HandleHttpError(w, err)

		return
	}

	if err := h.deleteReviewUseCase.Execute(user.ID, reviewID); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) listFeed(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromCtx(r.Context())
	if err != nil {
		httpErr := errors.NewHTTPErr("ocorreu um erro de autenticação, tente logar novamente", 401, "HANDLER:AUTHENTICATE_USER:INVALID_USER")
		errors.HandleHttpError(w, httpErr)

		return
	}

	strategy := r.URL.Query().Get("strategy")

	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.listFeedUsecase.Execute(user.ID, strategy, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := sendJsonResponse(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}
