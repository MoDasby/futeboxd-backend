package middleware

import (
	"context"
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

					next.ServeHTTP(w, r.WithContext(ctx))

					return
				}

				err := errors.NewHTTPErr(
					"sessão inválida",
					401,
					"MIDDLEWARE:AUTHENTICATION:INVALID_TOKEN",
				)
				errors.HandleHttpError(w, err)

				return
			}

			session, err := sessionRepository.FindOneByToken(r.Context(), token)
			if err != nil {
				newCookie := cookies.DeleteSessionCookie(cfg)

				http.SetCookie(w, newCookie)
				errors.HandleHttpError(w, err)

				return
			}

			if !session.IsValid() {
				if err := sessionRepository.Delete(r.Context(), session.ID); err != nil {
					errors.HandleHttpError(w, err)

					return
				}

				newCookie := cookies.DeleteSessionCookie(cfg)

				http.SetCookie(w, newCookie)

				err := errors.NewHTTPErr(
					"essa sessão está expirada",
					401,
					"MIDDLEWARE:AUTHENTICATION:EXPIRED_TOKEN",
				)
				errors.HandleHttpError(w, err)

				return
			}

			// renova o token se o tempo de expiração for menor que 3 dias
			if time.Until(session.ExpiresAt).Hours() < 72 {
				session.ExpiresAt = time.Now().Add(session.DefaultExpiration())

				if err := sessionRepository.Update(r.Context(), session); err != nil {
					errors.HandleHttpError(w, err)

					return
				}

				newCookie := cookies.CreateSessionCookie(cfg, session.Token, session.ExpiresAt)

				http.SetCookie(w, newCookie)
			}

			ctx := context.WithValue(r.Context(), SessionKey, session)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
