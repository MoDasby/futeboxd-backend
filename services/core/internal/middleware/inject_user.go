package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type ctxKey string

type Middleware func(next http.HandlerFunc) http.HandlerFunc

const (
	UserKey    ctxKey = "user"
	SessionKey ctxKey = "session"
)

func NewInjectUserMiddleware(
	sessionRepository domain.SessionRepository,
	userRepository domain.UserRepository,
) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				user := domain.NewAnonymousUser()

				ctx := context.WithValue(r.Context(), UserKey, user)

				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			session, err := sessionRepository.FindOneByToken(token)
			if err != nil {
				errors.HandleHttpError(w, err)

				return
			}

			if isSessionExpired(session) {
				err := errors.NewHTTPErr(
					"essa sessão está expirada",
					401,
					"MIDDLEWARE:AUTHENTICATION:EXPIRED_TOKEN",
				)
				errors.HandleHttpError(w, err)

				return
			}

			user, err := userRepository.FindOneByIdOrUsername(session.UserID)
			if err != nil {
				errors.HandleHttpError(w, err)

				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			ctx = context.WithValue(ctx, SessionKey, session)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func isSessionExpired(session *domain.Session) bool {
	return session.ExpiresAt.UTC().Before(time.Now().UTC())
}
