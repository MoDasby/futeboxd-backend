package reviews

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-backend/core/internal/review/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/auth"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type ReviewHandler struct {
	usecase ReviewUsecases
}

func NewReviewsHandler(usecase ReviewUsecases) *ReviewHandler {
	return &ReviewHandler{usecase: usecase}
}

func (h *ReviewHandler) RegisterRoutes(r *http.ServeMux, injectUser auth.Middleware) {
	r.HandleFunc("POST /reviews", injectUser(h.create, false))
	r.HandleFunc("GET /reviews/feed", injectUser(h.listFeed, false))
	r.HandleFunc("GET /reviews", injectUser(h.listAll, true))
	r.HandleFunc("DELETE /reviews/{reviewID}", injectUser(h.delete, false))
	r.HandleFunc("POST /reviews/{reviewID}/like", injectUser(h.toggleLikeReview, false))
}

func (h *ReviewHandler) create(w http.ResponseWriter, r *http.Request) {
	var input dto.ReviewInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpErr := errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:REVIEW:CREATE:INVALID_BODY",
		)

		errors.HandleHttpError(w, httpErr)

		return
	}

	if err := h.usecase.Create(r.Context(), &input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) listFeed(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.usecase.ListFeed(r.Context(), page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) listAll(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	options := ListReviewsOptions{
		Username:   r.URL.Query().Get("username"),
		Team:       r.URL.Query().Get("team"),
		Match:      r.URL.Query().Get("match"),
		SearchTerm: r.URL.Query().Get("search"),
	}

	output, err := h.usecase.ListBy(r.Context(), options, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) delete(w http.ResponseWriter, r *http.Request) {
	reviewID, err := utils.ParseIntValue(r.PathValue("reviewID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = errors.NewHTTPErr(
				httpErr.Msg+"para review_id",
				httpErr.Code,
				httpErr.Context+"INVALID_REVIEW_ID",
			)
		}

		errors.HandleHttpError(w, err)

		return
	}

	if err := h.usecase.Delete(r.Context(), reviewID); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *ReviewHandler) toggleLikeReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := utils.ParseIntValue(r.PathValue("reviewID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = errors.NewHTTPErr(
				httpErr.Msg+"para review_id",
				httpErr.Code,
				httpErr.Context+"INVALID_REVIEW_ID",
			)
		}

		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.usecase.ToggleLike(r.Context(), reviewID)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}
