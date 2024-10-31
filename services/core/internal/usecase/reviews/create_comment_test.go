package usecase

/*
import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateComment(t *testing.T) {
	mockUserRepo := new(repository.MockUserRepo)
	mockCommentsRepo := new(repository.MockCommentRepo)
	mockReviewRepo := new(repository.MockReviewRepo)

	uc := NewCreateCommentUsecase(mockCommentsRepo, mockUserRepo, mockReviewRepo)

	input := CommentInputDTO{
		AuthorID: "uuid",
		ParentID: 12,
		Content:  "Essa é a pior partida que eu já vi na vida",
	}

	user := domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	mockUserRepo.On("FindOneByIdOrUsername", "uuid").Return(&user, nil)
	mockCommentsRepo.On("Create", mock.AnythingOfType("*domain.Comment")).Return(nil)

	err := uc.Execute(input)

	assert.Nil(t, err)

	mockUserRepo.AssertExpectations(t)
	mockCommentsRepo.AssertExpectations(t)
} */
