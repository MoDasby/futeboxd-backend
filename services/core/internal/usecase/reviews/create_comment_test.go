package usecase

import (
	"testing"
	"time"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	mockCommentsRepo := repository_mocks.NewMockCommentsRepo()
	mockReviewRepo := repository_mocks.NewMockReviewRepo()

	uc := NewCreateCommentUsecase(mockCommentsRepo, mockReviewRepo)

	user := &domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	review := &domain.Review{
		ID:          1,
		Author:      user,
		Rate:        2,
		Description: "partida legal",
		MatchID:     19754,
		HomeTeamID:  7632,
		AwayTeamID:  7633,
		Likes:       0,
		IsLiked:     false,
		CreatedAt:   time.Now(),
	}

	mockReviewRepo.Create(review)

	input := CommentInputDTO{
		Author:   user,
		ParentID: 1,
		Content:  "Essa é a pior partida que eu já vi na vida",
	}

	err := uc.Execute(input)

	assert.Nil(t, err)

	input = CommentInputDTO{
		Author:   user,
		ParentID: 12,
		Content:  "Essa é a pior partida que eu já vi na vida",
	}

	err = uc.Execute(input)

	assert.NotNil(t, err)
}
