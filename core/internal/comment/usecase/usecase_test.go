package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/modasby/futeboxd-backend/core/internal/comment"
	"github.com/modasby/futeboxd-backend/core/internal/comment/dto"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/stretchr/testify/assert"

	commentMock "github.com/modasby/futeboxd-backend/core/internal/comment/mock"
	reviewMock "github.com/modasby/futeboxd-backend/core/internal/review/mock"
	userMock "github.com/modasby/futeboxd-backend/core/internal/user/mock"
)

func TestCommentUsecases_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks das dependências
	mockCommentRepo := commentMock.NewMockRepository(ctrl)
	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	// Instância do usecase com os mocks
	uc := NewCommentUsecases(mockCommentRepo, mockReviewRepo, mockUserRepo)

	// Dados de entrada
	input := &dto.CommentInput{
		ParentID: 1,
		Content:  "Great review!",
	}

	// Cenário 1: Sucesso ao criar um comentário
	t.Run("Success", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de review
		mockReviewRepo.EXPECT().
			ExistsByID(ctx, input.ParentID).
			Return(true, nil)

		// Mock do repositório de usuário
		author := &user.User{ID: session.UserID, Username: "testuser"}
		mockUserRepo.EXPECT().
			FindOneByIdOrUsername(ctx, session.UserID).
			Return(author, nil)

		// Mock do repositório de comentário
		mockCommentRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(nil)

		// Executa o método
		err := uc.Create(ctx, input)

		// Verifica o resultado
		assert.NoError(t, err)
	})

	// Cenário 2: Review não existe
	t.Run("ReviewNotFound", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de review
		mockReviewRepo.EXPECT().
			ExistsByID(ctx, input.ParentID).
			Return(false, nil)

		// Executa o método
		err := uc.Create(ctx, input)

		// Verifica o erro
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Review especificada não existe")
	})

	// Cenário 3: Erro ao buscar o autor
	t.Run("AuthorNotFound", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de review
		mockReviewRepo.EXPECT().
			ExistsByID(ctx, input.ParentID).
			Return(true, nil)

		// Mock do repositório de usuário
		mockUserRepo.EXPECT().
			FindOneByIdOrUsername(ctx, session.UserID).
			Return(nil, errors.New("user not found"))

		// Executa o método
		err := uc.Create(ctx, input)

		// Verifica o erro
		assert.Error(t, err)
	})
}

func TestCommentUsecases_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks das dependências
	mockCommentRepo := commentMock.NewMockRepository(ctrl)
	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	// Instância do usecase com os mocks
	uc := NewCommentUsecases(mockCommentRepo, mockReviewRepo, mockUserRepo)

	// Cenário 1: Sucesso ao deletar um comentário
	t.Run("Success", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de comentário
		comment := &comment.Comment{ID: 1, Author: &user.User{ID: session.UserID}}
		mockCommentRepo.EXPECT().
			FindOneByID(ctx, session.UserID, int64(1)).
			Return(comment, nil)

		mockCommentRepo.EXPECT().
			Delete(ctx, int64(1)).
			Return(nil)

		// Executa o método
		err := uc.Delete(ctx, 1)

		// Verifica o resultado
		assert.NoError(t, err)
	})

	// Cenário 2: Usuário não é o autor do comentário
	t.Run("Forbidden", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de comentário
		comment := &comment.Comment{ID: 1, Author: &user.User{ID: uuid.NewString()}}
		mockCommentRepo.EXPECT().
			FindOneByID(ctx, session.UserID, int64(1)).
			Return(comment, nil)

		// Executa o método
		err := uc.Delete(ctx, 1)

		// Verifica o erro
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Você não pode executar essa ação")
	})
}

func TestCommentUsecases_ListByReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks das dependências
	mockCommentRepo := commentMock.NewMockRepository(ctrl)
	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	// Instância do usecase com os mocks
	uc := NewCommentUsecases(mockCommentRepo, mockReviewRepo, mockUserRepo)

	// Cenário 1: Sucesso ao listar comentários
	t.Run("Success", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de comentário
		comments := []comment.Comment{
			{
				ID:        1,
				Author:    &user.User{ID: uuid.NewString(), Username: "testuser"},
				Content:   "Great review!",
				LikeCount: 10,
				IsLiked:   true,
				CreatedAt: time.Now(),
			},
		}
		mockCommentRepo.EXPECT().
			ListByReview(ctx, session.UserID, int64(1), gomock.Any()).
			Return(comments, nil)

		// Executa o método
		result, err := uc.ListByReview(ctx, 1, &pagination.Page{})

		// Verifica o resultado
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Great review!", result[0].Content)
	})
}

func TestCommentUsecases_ToggleLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks das dependências
	mockCommentRepo := commentMock.NewMockRepository(ctrl)
	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)

	// Instância do usecase com os mocks
	uc := NewCommentUsecases(mockCommentRepo, mockReviewRepo, mockUserRepo)

	// Cenário 1: Sucesso ao curtir um comentário
	t.Run("Like", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de comentário
		comment := &comment.Comment{ID: 1, IsLiked: false}
		mockCommentRepo.EXPECT().
			FindOneByID(ctx, session.UserID, int64(1)).
			Return(comment, nil)

		mockCommentRepo.EXPECT().
			Like(ctx, session.UserID, int64(1)).
			Return(nil)

		// Executa o método
		result, err := uc.ToggleLike(ctx, 1)

		// Verifica o resultado
		assert.NoError(t, err)
		assert.True(t, result.Like)
	})

	// Cenário 2: Sucesso ao descurtir um comentário
	t.Run("Unlike", func(t *testing.T) {
		session := &session.Session{UserID: uuid.NewString()}
		ctx := context.WithValue(context.TODO(), middleware.SessionKey, session)

		// Mock do repositório de comentário
		comment := &comment.Comment{ID: 1, IsLiked: true}
		mockCommentRepo.EXPECT().
			FindOneByID(ctx, session.UserID, int64(1)).
			Return(comment, nil)

		mockCommentRepo.EXPECT().
			Unlike(ctx, session.UserID, int64(1)).
			Return(nil)

		// Executa o método
		result, err := uc.ToggleLike(ctx, 1)

		// Verifica o resultado
		assert.NoError(t, err)
		assert.False(t, result.Like)
	})
}
