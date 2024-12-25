package usecase

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/comment"
	"github.com/modasby/futeboxd-backend/core/internal/comment/dto"
	"github.com/modasby/futeboxd-backend/core/internal/review"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type commentUsecases struct {
	commentRepo comment.Repository
	reviewRepo  review.Repository
	userRepo    user.Repository
}

func NewCommentUsecases(
	commentRepo comment.Repository,
	reviewRepo review.Repository,
	userRepo user.Repository,
) comment.CommentUsecases {
	return &commentUsecases{commentRepo: commentRepo, reviewRepo: reviewRepo, userRepo: userRepo}
}

func (uc *commentUsecases) Create(ctx context.Context, input *dto.CommentInput) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	reviewExists, err := uc.reviewRepo.ExistsByID(ctx, input.ParentID)
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

	author, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return err
	}

	comment, err := comment.NewComment(author, input.ParentID, input.Content)
	if err != nil {
		return err
	}

	if err := uc.commentRepo.Create(ctx, comment); err != nil {
		return err
	}

	return nil
}

func (uc *commentUsecases) Delete(ctx context.Context, commentID int64) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	comment, err := uc.commentRepo.FindOneByID(ctx, session.UserID, commentID)
	if err != nil {
		return err
	}

	if session.UserID != comment.Author.ID {
		return errors.NewHTTPErr(
			"você não pode executar essa ação",
			403,
			"USECASE:REVIEWS:DELETE_COMMENT:FORBIDDEN",
		)
	}

	if err := uc.commentRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	return nil
}

func (uc *commentUsecases) ListByReview(ctx context.Context, reviewID int64, page *pagination.Page) ([]dto.Comment, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	comments, err := uc.commentRepo.ListByReview(ctx, session.UserID, reviewID, page)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Comment, 0)

	for _, comment := range comments {
		comment := dto.Comment{
			ID: comment.ID,
			Author: dto.ContentAuthor{
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

func (uc *commentUsecases) ToggleLike(ctx context.Context, commentID int64) (*dto.LikeStats, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	comment, err := uc.commentRepo.FindOneByID(ctx, session.UserID, commentID)
	if err != nil {
		return nil, err
	}

	if comment.IsLiked {
		if err := uc.commentRepo.Unlike(ctx, session.UserID, commentID); err != nil {
			return nil, err
		}

		return &dto.LikeStats{Like: false}, nil
	}

	if err := uc.commentRepo.Like(ctx, session.UserID, commentID); err != nil {
		return nil, err
	}

	return &dto.LikeStats{Like: true}, nil
}
