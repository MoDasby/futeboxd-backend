package usecase

import (
	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type DeleteReviewUseCase struct {
	reviewRepository domain.ReviewRepository
}

func NewDeleteReviewUseCase(repo domain.ReviewRepository) *DeleteReviewUseCase {
	return &DeleteReviewUseCase{
		reviewRepository: repo,
	}
}

func (uc *DeleteReviewUseCase) Execute(reviewID int64, userID string) error {
	review, err := uc.reviewRepository.FindOneByID(reviewID)
	if err != nil {
		return err
	}

	if review.Author.ID != userID {
		return errors.NewHTTPErr(
			"você não pode executar essa ação",
			403,
			"USECASE:REVIEWS:DELETE_REVIEW:FORBIDDEN",
		)
	}

	if err := uc.reviewRepository.Delete(reviewID); err != nil {
		return err
	}

	return nil
}
