package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/comment"
	"github.com/modasby/futeboxd-backend/core/internal/comment/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type Handler struct {
	usecase comment.CommentUsecases
}

func NewCommentHandler(usecase comment.CommentUsecases) *Handler {
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
	reviewID, err := utils.ParseIntValue(r.Context(), r.PathValue("reviewID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = &errors.HTTPErr{
				Msg:        httpErr.Msg + "para review_id",
				Code:       httpErr.Code,
				Context:    httpErr.Context + "INVALID_REVIEW_ID",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  r.Context().Value("traceID").(string),
				Timestamp:  time.Now().UTC(),
			}
		}

		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	var input dto.CommentInput

	input.ParentID = reviewID

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpErr := &errors.HTTPErr{
			Msg:        "Corpo de requisição inválido",
			Code:       http.StatusBadRequest,
			Context:    "REVIEW:HANDLER:CREATE_COMMENT:INVALID_BODY",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  r.Context().Value("traceID").(string),
			Timestamp:  time.Now().UTC(),
		}

		errors.HandleHttpError(r.Context(), w, httpErr)

		return
	}

	if err := h.usecase.Create(r.Context(), &input); err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.ParseIntValue(r.Context(), r.PathValue("commentID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = &errors.HTTPErr{
				Msg:        httpErr.Msg + "para comment_id",
				Code:       httpErr.Code,
				Context:    httpErr.Context + "INVALID_COMMENT_ID",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  r.Context().Value("traceID").(string),
				Timestamp:  time.Now().UTC(),
			}
		}

		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	if err := h.usecase.Delete(r.Context(), commentID); err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *Handler) listComments(w http.ResponseWriter, r *http.Request) {
	reviewID, err := utils.ParseIntValue(r.Context(), r.URL.Query().Get("review"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = &errors.HTTPErr{
				Msg:        httpErr.Msg + "para review_id",
				Code:       httpErr.Code,
				Context:    httpErr.Context + "INVALID_REVIEW_ID",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  r.Context().Value("traceID").(string),
				Timestamp:  time.Now().UTC(),
			}
		}

		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	page, err := pagination.WithRequest(r)
	if err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	output, err := h.usecase.ListByReview(r.Context(), reviewID, page)
	if err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}
}

func (h *Handler) toggleLikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.ParseIntValue(r.Context(), r.PathValue("commentID"))
	if err != nil {
		httpErr, ok := err.(*errors.HTTPErr)

		if ok {
			err = &errors.HTTPErr{
				Msg:        httpErr.Msg + "para comment_id",
				Code:       httpErr.Code,
				Context:    httpErr.Context + "INVALID_COMMENT_ID",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  r.Context().Value("traceID").(string),
				Timestamp:  time.Now().UTC(),
			}
		}

		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	output, err := h.usecase.ToggleLike(r.Context(), commentID)
	if err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}
}
