package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/pkg/cookies"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type AuthMiddleware func(next http.HandlerFunc, permitAnonymous bool) http.HandlerFunc

type ctxKey string

const (
	SessionKey ctxKey = "session"
)

func NewAuthMiddleware(
	sessionRepository session.Repository, cfg *config.Cookies,
) AuthMiddleware {
	return func(next http.HandlerFunc, permitAnonymous bool) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, _ := r.Cookie("sid")
			tokenHeader := r.Header.Get("Authorization")

			var token string

			if tokenCookie != nil {
				token = tokenCookie.Value
			}

			if tokenHeader != "" {
				token = tokenHeader
			}

			if token == "" {
				if permitAnonymous {
					session := session.NewAnonymousSession()

					ctx := context.WithValue(r.Context(), SessionKey, session)

					slog.InfoContext(ctx, "permitindo acesso anônimo", "clientIP", r.RemoteAddr)

					next.ServeHTTP(w, r.WithContext(ctx))

					return
				}

				err := &errors.HTTPErr{
					Msg:        "Sessão inválida",
					Code:       http.StatusUnauthorized,
					Context:    "MIDDLEWARE:AUTHENTICATION:INVALID_TOKEN",
					StackTrace: errors.CaptureStackTrace(),
					ErrorCode:  r.Context().Value("traceID").(string),
					Timestamp:  time.Now().UTC(),
				}
				errors.HandleHttpError(r.Context(), w, err)

				return
			}

			session, err := sessionRepository.FindOneByToken(r.Context(), token)
			if err != nil {
				newCookie := cookies.DeleteSessionCookie(cfg)

				http.SetCookie(w, newCookie)
				errors.HandleHttpError(r.Context(), w, err)

				return
			}

			if !session.IsValid() {
				if err := sessionRepository.Delete(r.Context(), session.ID); err != nil {
					errors.HandleHttpError(r.Context(), w, err)

					return
				}

				newCookie := cookies.DeleteSessionCookie(cfg)

				http.SetCookie(w, newCookie)

				err := &errors.HTTPErr{
					Msg:        "Essa sessão está expirada",
					Code:       http.StatusUnauthorized,
					Context:    "MIDDLEWARE:AUTHENTICATION:EXPIRED_TOKEN",
					StackTrace: errors.CaptureStackTrace(),
					ErrorCode:  r.Context().Value("traceID").(string),
					Timestamp:  time.Now().UTC(),
				}
				errors.HandleHttpError(r.Context(), w, err)

				return
			}

			// renova o token se o tempo de expiração for menor que 3 dias
			if time.Until(session.ExpiresAt).Hours() < 72 {
				session.ExpiresAt = time.Now().Add(session.DefaultExpiration())

				if err := sessionRepository.Update(r.Context(), session); err != nil {
					errors.HandleHttpError(r.Context(), w, err)

					return
				}

				newCookie := cookies.CreateSessionCookie(cfg, session.Token, session.ExpiresAt)

				http.SetCookie(w, newCookie)
			}

			slog.Info(
				"Requisição autenticada",
				"sessionID", session.ID,
				"userID", session.UserID,
			)

			ctx := context.WithValue(r.Context(), SessionKey, session)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
