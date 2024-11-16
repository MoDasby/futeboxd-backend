package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type CreateReviewInputDTO struct {
	Author      *domain.User
	Rate        int    `json:"rate"`
	Description string `json:"description"`
	MatchID     int64  `json:"match_id"`
}

type CreateReviewUseCase struct {
	reviewRepository domain.ReviewRepository
	footballClient   football.Client
}

func NewCreateReviewUseCase(
	reviewRepository domain.ReviewRepository,
	footballClient football.Client,
) *CreateReviewUseCase {
	return &CreateReviewUseCase{
		reviewRepository: reviewRepository,
		footballClient:   footballClient,
	}
}

func (uc *CreateReviewUseCase) Execute(input CreateReviewInputDTO) error {
	match, err := uc.footballClient.GetMatch(input.MatchID)
	if err != nil {
		return err
	}

	review, err := domain.NewReview(
		input.Author,
		input.Rate,
		match.ID,
		input.Description,
		match.HomeCompetitor.Team.ID,
		match.AwayCompetitor.Team.ID,
	)
	if err != nil {
		return err
	}

	err = uc.reviewRepository.Create(review)
	if err != nil {
		return err
	}

	return nil
}
