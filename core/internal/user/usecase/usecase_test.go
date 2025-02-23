package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/internal/user/dto"
	userMock "github.com/modasby/futeboxd-backend/core/internal/user/mock"
	emailMock "github.com/modasby/futeboxd-backend/core/pkg/email/mock"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	footballMock "github.com/modasby/futeboxd-backend/core/pkg/football/mock"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.UserInput{
		Username:       "testuser",
		Name:           "Test User",
		Email:          "testuser@example.com",
		Password:       "password123",
		FavoriteTeamID: 1,
	}

	mockRepo.EXPECT().Exists(context.Background(), input.Username, input.Email).Return(false, nil)
	footballClientMock.EXPECT().GetTeam(context.Background(), input.FavoriteTeamID).Return(nil, nil)
	mockRepo.EXPECT().Create(context.Background(), gomock.Any()).Return(nil)

	err := uc.CreateUser(context.Background(), input)
	assert.NoError(t, err)
}

func TestCreateUser_DuplicateUsernameOrEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.UserInput{
		Username:       "testuser",
		Name:           "Test User",
		Email:          "testuser@example.com",
		Password:       "password123",
		FavoriteTeamID: 1,
	}

	// Expectation for checking if user exists
	mockRepo.EXPECT().Exists(context.Background(), input.Username, input.Email).Return(true, nil)

	err := uc.CreateUser(context.Background(), input)
	assert.Error(t, err)
}

func TestCreateUser_FavoriteTeamNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.UserInput{
		Username:       "testuser",
		Name:           "Test User",
		Email:          "testuser@example.com",
		Password:       "password123",
		FavoriteTeamID: 1,
	}

	// Expectation for checking if user exists
	mockRepo.EXPECT().Exists(context.Background(), input.Username, input.Email).Return(false, nil)
	// Expectation for fetching the favorite team
	footballClientMock.EXPECT().GetTeam(context.Background(), input.FavoriteTeamID).Return(nil, errors.New("team not found"))

	err := uc.CreateUser(context.Background(), input)
	assert.Error(t, err)
	assert.EqualError(t, err, "team not found")
}

func TestCreateUser_ErrorOnCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.UserInput{
		Username:       "testuser",
		Name:           "Test User",
		Email:          "testuser@example.com",
		Password:       "password123",
		FavoriteTeamID: 1,
	}

	// Expectation for checking if user exists
	mockRepo.EXPECT().Exists(context.Background(), input.Username, input.Email).Return(false, nil)
	// Expectation for fetching the favorite team
	footballClientMock.EXPECT().GetTeam(context.Background(), input.FavoriteTeamID).Return(nil, nil)
	// Expectation for creating the user
	mockRepo.EXPECT().Create(context.Background(), gomock.Any()).Return(errors.New("failed to create user"))

	err := uc.CreateUser(context.Background(), input)
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to create user")
}

func TestEditUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.EditUser{
		Username:       "newusername",
		Email:          "newemail@example.com",
		Name:           "New Name",
		Bio:            "New bio",
		FavoriteTeamID: 2,
	}

	// Mock session data
	session := &session.Session{
		UserID: "user123",
	}
	ctx := context.WithValue(context.Background(), middleware.SessionKey, session)

	// Expectation for fetching the user
	mockUser := &user.User{
		ID:       "user123",
		Email:    "oldemail@example.com",
		Username: "oldusername",
		Password: "12345678",
	}
	mockRepo.EXPECT().FindOneByIdOrUsername(ctx, session.UserID).Return(mockUser, nil)

	// Expectation for checking if email exists
	mockRepo.EXPECT().Exists(ctx, "", input.Email).Return(false, nil)

	// Expectation for checking if username exists
	mockRepo.EXPECT().Exists(ctx, input.Username, "").Return(false, nil)

	// Expectation for fetching the team
	footballClientMock.EXPECT().GetTeam(gomock.Any(), input.FavoriteTeamID).Return(&football.Team{ID: input.FavoriteTeamID}, nil)

	// Expectation for updating the user
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := uc.EditUser(ctx, input)
	assert.NoError(t, err)
}

func TestEditUser_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.EditUser{
		Email: "newemail@example.com",
	}

	// Mock session data
	session := &session.Session{
		UserID: "user123",
	}
	ctx := context.WithValue(context.Background(), middleware.SessionKey, session)

	// Mock the user retrieval
	mockUser := &user.User{
		ID:       "user123",
		Email:    "oldemail@example.com",
		Username: "oldusername",
	}
	mockRepo.EXPECT().FindOneByIdOrUsername(ctx, session.UserID).Return(mockUser, nil)

	// Expectation for checking if email exists
	mockRepo.EXPECT().Exists(ctx, "", input.Email).Return(true, nil)

	err := uc.EditUser(ctx, input)
	assert.Error(t, err)
}

func TestEditUser_UsernameAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.EditUser{
		Username: "newusername",
	}

	// Mock session data
	session := &session.Session{
		UserID: "user123",
	}
	ctx := context.WithValue(context.Background(), middleware.SessionKey, session)

	// Mock the user retrieval
	mockUser := &user.User{
		ID:       "user123",
		Email:    "oldemail@example.com",
		Username: "oldusername",
	}
	mockRepo.EXPECT().FindOneByIdOrUsername(ctx, session.UserID).Return(mockUser, nil)

	// Expectation for checking if username exists
	mockRepo.EXPECT().Exists(ctx, input.Username, "").Return(true, nil)

	err := uc.EditUser(ctx, input)
	assert.Error(t, err)
}

func TestEditUser_TeamNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userMock.NewMockRepository(ctrl)
	mockEmailClient := emailMock.NewMockClient(ctrl)
	footballClientMock := footballMock.NewMockClient(ctrl)

	uc := NewUsersUsecases(mockRepo, footballClientMock, mockEmailClient, config.Frontend{})

	input := &dto.EditUser{
		FavoriteTeamID: 2,
	}

	// Mock session data
	session := &session.Session{
		UserID: "user123",
	}
	ctx := context.WithValue(context.Background(), middleware.SessionKey, session)

	// Mock the user retrieval
	mockUser := &user.User{
		ID:       "user123",
		Email:    "oldemail@example.com",
		Username: "oldusername",
	}
	mockRepo.EXPECT().FindOneByIdOrUsername(ctx, session.UserID).Return(mockUser, nil)

	// Expectation for fetching the team
	footballClientMock.EXPECT().GetTeam(gomock.Any(), input.FavoriteTeamID).Return(nil, errors.New("team not found"))

	err := uc.EditUser(ctx, input)
	assert.Error(t, err)
	assert.EqualError(t, err, "team not found")
}
