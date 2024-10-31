package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type ctxKey string

type Middleware func(next http.HandlerFunc, permitAnonymous bool) http.HandlerFunc

const (
	UserKey    ctxKey = "user"
	SessionKey ctxKey = "session"
)

func NewInjectUserMiddleware(
	sessionRepository domain.SessionRepository,
	userRepository domain.UserRepository,
) Middleware {
	return func(next http.HandlerFunc, permitAnonymous bool) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				if permitAnonymous {
					user := domain.NewAnonymousUser()

					ctx := context.WithValue(r.Context(), UserKey, user)

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

			session, err := sessionRepository.FindOneByToken(token)
			if err != nil {
				errors.HandleHttpError(w, err)

				return
			}

			if !session.IsValid() {
				sessionRepository.Delete(session.ID)

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

				sessionRepository.Update(session)
			}

			// TODO colocar o user como uma entidade dentro do seession
			// pra nao precisar fazer duas queries
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
