package usecase

import "github.com/modasby/futeboxd-api/services/core/internal/domain"

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

func (uc GetCommentsUsecase) Execute()
