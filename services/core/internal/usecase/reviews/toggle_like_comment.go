package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
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

type outputDTO struct {
	Like bool `json:"like"`
}

func (uc *ToggleLikeCommentUsecase) Execute(requesterID string, commentID int64) (*outputDTO, error) {
	commentExists, err := uc.commentRepo.ExistsByID(commentID)
	if err != nil {
		return nil, err
	}

	if !commentExists {
		return nil, errors.NewHTTPErr(
			"esse comentário não existe",
			404,
			"USECASE:TOGGLE_LIKE_COMMENT:COMMENT_NOT_FOUND",
		)
	}

	isLiked, err := uc.commentRepo.IsLiked(requesterID, commentID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		if err := uc.commentRepo.Unlike(requesterID, commentID); err != nil {
			return nil, err
		}

		return &outputDTO{Like: false}, nil
	}

	if err := uc.commentRepo.Like(requesterID, commentID); err != nil {
		return nil, err
	}

	return &outputDTO{Like: true}, nil
}
