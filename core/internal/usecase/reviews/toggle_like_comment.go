package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type ToggleLikeCommentUsecase struct {
	commentRepo domain.CommentsRepository
}

func NewToggleLikeCommentUsecase(
	commentsRepo domain.CommentsRepository,
) *ToggleLikeCommentUsecase {
	return &ToggleLikeCommentUsecase{
		commentRepo: commentsRepo,
	}
}

type OutputDTO struct {
	Like bool `json:"like"`
}

func (uc *ToggleLikeCommentUsecase) Execute(requesterID string, commentID int64) (*OutputDTO, error) {
	comment, err := uc.commentRepo.FindOneByID(requesterID, commentID)
	if err != nil {
		return nil, err
	}

	if comment.IsLiked {
		if err := uc.commentRepo.Unlike(requesterID, commentID); err != nil {
			return nil, err
		}

		return &OutputDTO{Like: false}, nil
	}

	if err := uc.commentRepo.Like(requesterID, commentID); err != nil {
		return nil, err
	}

	return &OutputDTO{Like: true}, nil
}
