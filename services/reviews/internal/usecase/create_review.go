package usecase

import (
	"fmt"

	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/pkg/client/user"
	"github.com/modasby/futeboxd-api/services/reviews/internal/domain"
)

type CreateReviewInputDTO struct {
	UserID         string `json:"user_id"`
	AuthorUsername string `json:"author_username"`
	Rate           int    `json:"rate"`
	Description    string `json:"description"`
	MatchID        int64  `json:"match_id"`
}

type CreateReviewUseCase struct {
	reviewRepository domain.ReviewRepository
	footballClient   football.Client
}

func NewCreateReviewUseCase(
	reviewRepository domain.ReviewRepository,
	footballClient football.Client,
	userClient user.Client,
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

	userID := input.UserID

	fmt.Println(match.Competitors[0].Team.ID)
	fmt.Println(match.Competitors[1].Team.ID)

	homeTeamID := match.Competitors[0].Team.ID
	awayTeamID := match.Competitors[1].Team.ID

	review, err := domain.NewReview(userID, input.Rate, match.ID, input.Description, homeTeamID, awayTeamID)
	if err != nil {
		return err
	}

	err = uc.reviewRepository.AddReview(review)
	if err != nil {
		return err
	}

	return nil
}
