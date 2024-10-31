package usecase

import (
	"github.com/modasby/futeboxd-api/pkg/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
)

type GetCommentsUsecase struct {
	commentsRepository domain.CommentsRepository
}

func NewGetCommentsUsecase(
	commentsRepository domain.CommentsRepository,
) *GetCommentsUsecase {
	return &GetCommentsUsecase{
		commentsRepository: commentsRepository,
	}
}

func (uc *GetCommentsUsecase) Execute(requesterID string, reviewID int64, page *pagination.Page) ([]dto.Comment, error) {
	comments, err := uc.commentsRepository.ListByReview(requesterID, reviewID, page)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Comment, 0)

	for _, comment := range comments {
		comment := dto.Comment{
			ID: comment.ID,
			Author: dto.UserDTO{
				ID:       comment.Author.ID,
				Username: comment.Author.Username,
			},
			Content:   comment.Content,
			LikeCount: comment.LikeCount,
			IsLiked:   comment.IsLiked,
			CreatedAt: comment.CreatedAt,
		}

		output = append(output, comment)
	}

	return output, nil
}
