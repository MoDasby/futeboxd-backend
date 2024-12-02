package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type DeleteCommentUsecase struct {
	commentsRepo domain.CommentsRepository
}

func NewDeleteCommentUsecase(commentsRepo domain.CommentsRepository) *DeleteCommentUsecase {
	return &DeleteCommentUsecase{
		commentsRepo: commentsRepo,
	}
}

func (uc *DeleteCommentUsecase) Execute(requesterID string, commentID int64) error {
	comment, err := uc.commentsRepo.FindOneByID(requesterID, commentID)
	if err != nil {
		return err
	}

	if requesterID != comment.Author.ID {
		return errors.NewHTTPErr(
			"você não pode executar essa ação",
			403,
			"USECASE:REVIEWS:DELETE_COMMENT:FORBIDDEN",
		)
	}

	if err := uc.commentsRepo.Delete(commentID); err != nil {
		return err
	}

	return nil
}
