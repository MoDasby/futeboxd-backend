package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/modasby/futeboxd-backend/core/internal/review"
	"github.com/modasby/futeboxd-backend/core/internal/review/dto"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/stretchr/testify/assert"

	reviewMock "github.com/modasby/futeboxd-backend/core/internal/review/mock"
	userMock "github.com/modasby/futeboxd-backend/core/internal/user/mock"
	footballClientMock "github.com/modasby/futeboxd-backend/core/pkg/football/mock"
)

func TestCreateReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewReviewUsecases(mockReviewRepo, mockFootballClient, mockUserRepo)

	session := &session.Session{UserID: uuid.NewString()}
	ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)
	input := &dto.ReviewInput{
		MatchID:     1,
		Rate:        5,
		Description: "Great match!",
	}

	mockFootballClient.EXPECT().GetMatch(ctx, input.MatchID).Return(&football.Match{
		ID:             1,
		HomeCompetitor: football.Competitor{Team: football.Team{ID: 7632}},
		AwayCompetitor: football.Competitor{Team: football.Team{ID: 7632}},
	}, nil)

	mockUserRepo.EXPECT().FindOneByIdOrUsername(ctx, session.UserID).Return(&user.User{ID: session.UserID, Username: "test_user"}, nil)

	mockReviewRepo.EXPECT().Upsert(ctx, gomock.Any()).Return(nil)

	err := uc.Create(ctx, input)
	assert.NoError(t, err)
}

func TestCreateReview_FailureInvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewReviewUsecases(mockReviewRepo, mockFootballClient, mockUserRepo)

	ctx := context.Background()
	input := &dto.ReviewInput{
		MatchID:     1,
		Rate:        5,
		Description: "Great match!",
	}

	err := uc.Create(ctx, input)
	assert.Error(t, err)
}

func TestDeleteReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewReviewUsecases(mockReviewRepo, mockFootballClient, mockUserRepo)

	session := &session.Session{UserID: uuid.NewString()}
	ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

	mockReviewRepo.EXPECT().FindOneByID(ctx, session.UserID, int64(1)).Return(&review.Review{
		Author: &user.User{ID: session.UserID},
	}, nil)
	mockReviewRepo.EXPECT().Delete(ctx, int64(1)).Return(nil)

	err := uc.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestDeleteReview_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewReviewUsecases(mockReviewRepo, mockFootballClient, mockUserRepo)

	session := &session.Session{UserID: uuid.NewString()}
	ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

	mockReviewRepo.EXPECT().FindOneByID(ctx, session.UserID, int64(1)).Return(&review.Review{
		Author: &user.User{ID: "other_user"},
	}, nil)

	err := uc.Delete(ctx, 1)
	assert.Error(t, err)
	assert.Equal(t, 403, err.(*errors.HTTPErr).Code)
}
