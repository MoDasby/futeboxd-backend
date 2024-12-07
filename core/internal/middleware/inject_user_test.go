package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInjectUser(t *testing.T) {
	mockSessionRepo := repository_mocks.NewMockSessionRepo()
	mockUserRepo := repository_mocks.NewMockUserRepo()

	injectUser := NewInjectUserMiddleware(mockSessionRepo, mockUserRepo)

	t.Run("should allow anonymous access when permitAnonymous is true", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		next := func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(UserKey).(*domain.User)
			assert.NotNil(t, user)
			assert.Equal(t, "anonymous", user.Username)
			assert.NoError(t, uuid.Validate(user.ID))
		}

		middleware := injectUser(next, true)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return error when token is missing and permitAnonymous is false", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		middleware := injectUser(func(w http.ResponseWriter, r *http.Request) {}, false)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "invalid_token")
		rr := httptest.NewRecorder()

		middleware := injectUser(func(w http.ResponseWriter, r *http.Request) {}, false)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error for expired session", func(t *testing.T) {
		session, _ := domain.NewSession("uuid")
		session.ExpiresAt = time.Now().Add(-1 * time.Hour)

		mockSessionRepo.Create(session)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", session.Token)
		rr := httptest.NewRecorder()

		middleware := injectUser(func(w http.ResponseWriter, r *http.Request) {}, false)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockSessionRepo.Delete(session.ID)
	})

	t.Run("should pass valid session and user to next handler", func(t *testing.T) {
		session, _ := domain.NewSession("uuid")

		mockSessionRepo.Create(session)

		user := &domain.User{ID: "uuid", Username: "test_user"}

		mockUserRepo.Create(user)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", session.Token)
		rr := httptest.NewRecorder()

		next := func(w http.ResponseWriter, r *http.Request) {
			u := r.Context().Value(UserKey).(*domain.User)
			s := r.Context().Value(SessionKey).(*domain.Session)

			assert.Equal(t, user, u)
			assert.Equal(t, session, s)
		}

		middleware := injectUser(next, false)
		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
