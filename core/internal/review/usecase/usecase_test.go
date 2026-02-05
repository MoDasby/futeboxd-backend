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
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
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

func TestListFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewReviewUsecases(mockReviewRepo, mockFootballClient, mockUserRepo)

	session := &session.Session{UserID: uuid.NewString()}
	ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)
	page := &pagination.Page{Size: 20, Index: 1}

	reviews := []review.Review{
		{
			ID:            1,
			Author:        &user.User{ID: "author1", Username: "user1"},
			Rate:          5,
			Description:   "Great match!",
			MatchID:       100,
			Likes:         10,
			CommentsCount: 3,
			IsLiked:       true,
		},
		{
			ID:            2,
			Author:        &user.User{ID: "author2", Username: "user2"},
			Rate:          3,
			MatchID:       200,
			Likes:         0,
			CommentsCount: 0,
			IsLiked:       false,
		},
	}

	mockReviewRepo.EXPECT().ListFeed(ctx, session.UserID, page).Return(reviews, nil)

	matchesMap := map[int64]football.Match{
		100: {ID: 100, HomeCompetitor: football.Competitor{Team: football.Team{ID: 1}}, AwayCompetitor: football.Competitor{Team: football.Team{ID: 2}}},
		200: {ID: 200, HomeCompetitor: football.Competitor{Team: football.Team{ID: 3}}, AwayCompetitor: football.Competitor{Team: football.Team{ID: 4}}},
	}
	mockFootballClient.EXPECT().GetMatchesMap(ctx, gomock.Any()).Return(matchesMap, nil)

	result, err := uc.ListFeed(ctx, page)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, 2, result[1].ID)
	assert.Equal(t, 10, result[0].Likes)
	assert.Equal(t, true, result[0].IsLiked)
}

func TestListFeed_FailureInvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	uc := NewReviewUsecases(mockReviewRepo, mockFootballClient, mockUserRepo)

	ctx := context.Background()
	page := &pagination.Page{Size: 20, Index: 1}

	result, err := uc.ListFeed(ctx, page)
	assert.Error(t, err)
	assert.Nil(t, result)
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
