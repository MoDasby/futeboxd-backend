package usecase

import (
	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type ToggleLikeReviewUsecase struct {
	reviewRepo domain.ReviewRepository
}

func NewToggleLikeReviewUsecase(
	reviewRepo domain.ReviewRepository,
) *ToggleLikeReviewUsecase {
	return &ToggleLikeReviewUsecase{
		reviewRepo: reviewRepo,
	}
}

func (uc *ToggleLikeReviewUsecase) Execute(requesterID string, reviewID int64) (*outputDTO, error) {
	reviewExists, err := uc.reviewRepo.ExistsByID(reviewID)
	if err != nil {
		return nil, err
	}

	if !reviewExists {
		return nil, errors.NewHTTPErr(
			"essa review não existe",
			404,
			"USECASE:TOGGLE_LIKE_REVIEW:REVIEW_NOT_FOUND",
		)
	}

	isLiked, err := uc.reviewRepo.IsLiked(requesterID, reviewID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		if err := uc.reviewRepo.Unlike(requesterID, reviewID); err != nil {
			return nil, err
		}

		return &outputDTO{Like: false}, nil
	}

	if err := uc.reviewRepo.Like(requesterID, reviewID); err != nil {
		return nil, err
	}

	return &outputDTO{Like: true}, nil
}
