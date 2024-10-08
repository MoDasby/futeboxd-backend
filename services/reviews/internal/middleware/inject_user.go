package middleware

import (
	"context"
	"net/http"

	"github.com/modasby/futeboxd-api/pkg/client/user"
	"github.com/modasby/futeboxd-api/pkg/errors"
)

type sessionKey int

const SessionKey sessionKey = 0

func NewInjectUserMiddleware(userClient user.Client) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			session, err := userClient.GetAuthenticatedUser(token)
			if err != nil {
				errors.HandleHttpError(w, err)

				return
			}

			ctx := context.WithValue(r.Context(), SessionKey, session)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
