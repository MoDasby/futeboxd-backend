package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/modasby/futeboxd-backend/core/internal/profile"
	profileMock "github.com/modasby/futeboxd-backend/core/internal/profile/mock"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	footballClientMock "github.com/modasby/futeboxd-backend/core/pkg/football/mock"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

func TestFindByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := profileMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)

	uc := NewProfileUsecases(mockRepo, mockFootballClient)

	ctx := context.WithValue(context.Background(), middleware.SessionKey, &session.Session{UserID: "123"})

	t.Run("success case with favorite team", func(t *testing.T) {
		username := "john_doe"

		mockRepo.EXPECT().FindOneByUsername(ctx, username, "123").Return(&profile.Profile{
			UserID:         "123",
			Name:           "John Doe",
			Bio:            "A bio",
			Username:       "john_doe",
			ProfilePicture: "profile.jpg",
			FavoriteTeam:   42,
			FollowersCount: 10,
			FollowingCount: 5,
			IsFollowing:    false,
			ReviewsCount:   3,
		}, nil)

		mockFootballClient.EXPECT().GetTeam(int64(42)).Return(&football.Team{ID: 42, Name: "Team A"}, nil)

		result, err := uc.FindByUsername(ctx, username)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "John Doe", result.Name.String())
		assert.Equal(t, "Team A", result.FavoriteTeam.Name)
	})

	t.Run("success case without favorite team", func(t *testing.T) {
		username := "john_doe"

		mockRepo.EXPECT().FindOneByUsername(ctx, username, "123").Return(&profile.Profile{
			UserID:         "123",
			Name:           "John Doe",
			Bio:            "A bio",
			Username:       "john_doe",
			ProfilePicture: "profile.jpg",
			FavoriteTeam:   0,
			FollowersCount: 10,
			FollowingCount: 5,
			IsFollowing:    false,
			ReviewsCount:   3,
		}, nil)

		result, err := uc.FindByUsername(ctx, username)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.FavoriteTeam)
	})

	t.Run("error case - user not found", func(t *testing.T) {
		username := "not_found"

		mockRepo.EXPECT().FindOneByUsername(ctx, username, "123").Return(nil, errors.New("user not found"))

		result, err := uc.FindByUsername(ctx, username)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestSearchByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := profileMock.NewMockRepository(ctrl)
	mockFootballClient := footballClientMock.NewMockClient(ctrl)

	uc := NewProfileUsecases(mockRepo, mockFootballClient)

	ctx := context.WithValue(context.Background(), middleware.SessionKey, &session.Session{UserID: "123"})

	t.Run("success case", func(t *testing.T) {
		page := &pagination.Page{Size: 10, Index: 1}
		username := "john"

		mockRepo.EXPECT().Search(ctx, "123", username, page).Return([]profile.Profile{
			{
				UserID:       "123",
				Name:         "John Doe",
				Bio:          "A bio",
				Username:     "john_doe",
				FavoriteTeam: 42,
			},
		}, nil)

		mockFootballClient.EXPECT().GetTeam(int64(42)).Return(&football.Team{ID: 42, Name: "Team A"}, nil)

		result, err := uc.SearchByUsername(ctx, username, page)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "John Doe", result[0].Name.String())
		assert.Equal(t, "Team A", result[0].FavoriteTeam.Name)
	})

	t.Run("empty username", func(t *testing.T) {
		result, err := uc.SearchByUsername(ctx, "", &pagination.Page{Size: 10, Index: 1})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestToggleFollow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := profileMock.NewMockRepository(ctrl)

	uc := NewProfileUsecases(mockRepo, nil)

	ctx := context.WithValue(context.Background(), middleware.SessionKey, &session.Session{UserID: "123"})

	t.Run("success case - follow user", func(t *testing.T) {
		username := "john_doe"

		mockRepo.EXPECT().FindOneByUsername(ctx, username, "123").Return(&profile.Profile{
			UserID:      "456",
			IsFollowing: false,
		}, nil)

		mockRepo.EXPECT().Follow(ctx, "123", "456").Return(nil)

		result, err := uc.ToggleFollow(ctx, username)

		assert.NoError(t, err)
		assert.True(t, result.Following)
	})

	t.Run("success case - unfollow user", func(t *testing.T) {
		username := "john_doe"

		mockRepo.EXPECT().FindOneByUsername(ctx, username, "123").Return(&profile.Profile{
			UserID:      "456",
			IsFollowing: true,
		}, nil)

		mockRepo.EXPECT().Unfollow(ctx, "123", "456").Return(nil)

		result, err := uc.ToggleFollow(ctx, username)

		assert.NoError(t, err)
		assert.False(t, result.Following)
	})

	t.Run("error case - user not found", func(t *testing.T) {
		username := "not_found"

		mockRepo.EXPECT().FindOneByUsername(ctx, username, "123").Return(nil, errors.New("user not found"))

		result, err := uc.ToggleFollow(ctx, username)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
