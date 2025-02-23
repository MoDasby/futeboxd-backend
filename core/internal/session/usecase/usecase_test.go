package usecase

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/session/dto"
	sessionMock "github.com/modasby/futeboxd-backend/core/internal/session/mock"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	userMock "github.com/modasby/futeboxd-backend/core/internal/user/mock"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := sessionMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewSessionUsecases(mockSessionRepo, mockUserRepo)

	input := &dto.Login{
		Credential: "test_user",
		Password:   "password123",
	}

	mockUser, _ := user.NewUser(context.Background(), "test_user", "", "", "user@email.com", "password123", 0)

	mockUserRepo.EXPECT().FindOneByCredential(context.Background(), input.Credential).Return(mockUser, nil)

	mockSession := &session.Session{
		ID:        uuid.NewString(),
		UserID:    mockUser.ID,
		Token:     uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(session.DEFAULT_EXPIRATION),
	}
	mockSessionRepo.EXPECT().Create(context.Background(), gomock.Any()).Return(mockSession, nil)

	output, err := uc.Login(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, mockSession.Token, output.Token)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := sessionMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewSessionUsecases(mockSessionRepo, mockUserRepo)

	input := &dto.Login{
		Credential: "invalid_user",
		Password:   "wrong_password",
	}

	mockUserRepo.EXPECT().FindOneByCredential(context.Background(), input.Credential).Return(nil, &errors.HTTPErr{
		Msg:        "Usuário ou senha inválidos",
		Code:       http.StatusBadRequest,
		Context:    "SESSION:USECASE:LOGIN:USER_NOT_FOUND",
		StackTrace: errors.CaptureStackTrace(),
		ErrorCode:  "trace",
		Timestamp:  time.Now().UTC(),
		Original:   nil,
	})

	output, err := uc.Login(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestLogout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := sessionMock.NewMockRepository(ctrl)
	uc := NewSessionUsecases(mockSessionRepo, nil)

	session := &session.Session{
		ID: uuid.NewString(),
	}
	ctx := context.WithValue(context.Background(), middleware.SessionKey, session)

	mockSessionRepo.EXPECT().Delete(ctx, session.ID).Return(nil)

	err := uc.Logout(ctx)
	assert.NoError(t, err)
}

func TestLogout_InvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := sessionMock.NewMockRepository(ctrl)
	uc := NewSessionUsecases(mockSessionRepo, nil)

	ctx := context.Background()

	err := uc.Logout(ctx)
	assert.Error(t, err)
}
