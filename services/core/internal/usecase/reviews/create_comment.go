package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type CreateCommentUsecase struct {
	commentsRepository domain.CommentsRepository
	reviewRepository   domain.ReviewRepository
}

func NewCreateCommentUsecase(
	commentsRepository domain.CommentsRepository,
	reviewRepository domain.ReviewRepository,
) *CreateCommentUsecase {
	return &CreateCommentUsecase{
		commentsRepository: commentsRepository,
		reviewRepository:   reviewRepository,
	}
}

type CommentInputDTO struct {
	Author   *domain.User
	ParentID int64  `json:"parent_id"`
	Content  string `json:"content"`
}

func (uc *CreateCommentUsecase) Execute(input CommentInputDTO) error {
	reviewExists, err := uc.reviewRepository.ExistsByID(input.ParentID)
	if err != nil {
		return err
	}

	if !reviewExists {
		return errors.NewHTTPErr(
			"review especificada não existe",
			400,
			"USECASE:CREATE_COMMENT:REVIEW_NOT_FOUND",
		)
	}

	comment, err := domain.NewComment(input.Author, input.ParentID, input.Content)
	if err != nil {
		return err
	}

	if err := uc.commentsRepository.Create(comment); err != nil {
		return err
	}

	return nil
}
