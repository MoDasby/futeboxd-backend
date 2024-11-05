package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type DeleteReviewUseCase struct {
	reviewRepository domain.ReviewRepository
}

func NewDeleteReviewUseCase(repo domain.ReviewRepository) *DeleteReviewUseCase {
	return &DeleteReviewUseCase{
		reviewRepository: repo,
	}
}

func (uc *DeleteReviewUseCase) Execute(requesterID string, reviewID int64) error {
	review, err := uc.reviewRepository.FindOneByID(requesterID, reviewID)
	if err != nil {
		return err
	}

	if review.Author.ID != requesterID {
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
