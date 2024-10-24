package usecase

import "github.com/modasby/futeboxd-api/services/core/internal/domain"

type CreateCommentUsecase struct {
	commentsRepository domain.CommentsRepository
	userRepository     domain.UserRepository
}

func NewCreateCommentUsecase(
	commentsRepository domain.CommentsRepository,
	userRepository domain.UserRepository,
) *CreateCommentUsecase {
	return &CreateCommentUsecase{
		commentsRepository: commentsRepository,
		userRepository:     userRepository,
	}
}

type CommentInputDTO struct {
	AuthorID string `json:"author_id"`
	ParentID int64  `json:"parent_id"`
	Content  string `json:"content"`
}

func (uc *CreateCommentUsecase) Execute(input CommentInputDTO) error {
	user, err := uc.userRepository.FindOneByIdOrUsername(input.AuthorID)
	if err != nil {
		return err
	}

	comment := domain.NewComment(*user, input.ParentID, input.Content)

	if err := uc.commentsRepository.Create(comment); err != nil {
		return err
	}

	return nil
}
