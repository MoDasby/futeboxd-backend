package comment

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-backend/core/internal/comment/dto"
	"github.com/modasby/futeboxd-backend/core/internal/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type Handler struct {
	usecase CommentUsecases
}

func NewCommentHandler(usecase CommentUsecases) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux, injectUser middleware.AuthMiddleware) {
	r.HandleFunc("POST /comments/{reviewID}", injectUser(h.createComment, false))
	r.HandleFunc("DELETE /comments/{commentID}", injectUser(h.deleteComment, false))
	r.HandleFunc("GET /comments", injectUser(h.listComments, true))
	r.HandleFunc("POST /comments/{commentID}/like", injectUser(h.toggleLikeComment, false))
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
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

	var input dto.CommentInput

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

	if err := h.usecase.Create(r.Context(), &input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.ParseIntValue(r.PathValue("commentID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = errors.NewHTTPErr(
				httpErr.Msg+"para comment_id",
				httpErr.Code,
				httpErr.Context+"INVALID_COMMENT_ID",
			)
		}

		errors.HandleHttpError(w, err)

		return
	}

	if err := h.usecase.Delete(r.Context(), commentID); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *Handler) listComments(w http.ResponseWriter, r *http.Request) {
	reviewID, err := utils.ParseIntValue(r.URL.Query().Get("review"))
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

	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.usecase.ListByReview(r.Context(), reviewID, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *Handler) toggleLikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.ParseIntValue(r.PathValue("commentID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = errors.NewHTTPErr(
				httpErr.Msg+"para comment_id",
				httpErr.Code,
				httpErr.Context+"INVALID_COMMENT_ID",
			)
		}

		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.usecase.ToggleLike(r.Context(), commentID)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}
