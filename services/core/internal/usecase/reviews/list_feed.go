package usecase

import (
	"fmt"

	"github.com/modasby/futeboxd-api/pkg/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
)

type ListFeedUsecase struct {
	reviewsRepo    domain.ReviewRepository
	userRepository domain.UserRepository
}

func NewListFeedUsecase(
	reviewsRepo domain.ReviewRepository,
	userRepository domain.UserRepository,
) *ListFeedUsecase {
	return &ListFeedUsecase{
		reviewsRepo:    reviewsRepo,
		userRepository: userRepository,
	}
}

func (uc *ListFeedUsecase) Execute(requesterID, strategy string, page *pagination.Page) ([]ReviewOutputDTO, error) {
	if strategy == "" {
		strategy = "relevant"
	}

	reviews, err := uc.reviewsRepo.ListFeed(requesterID, strategy, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	output := make([]ReviewOutputDTO, 0)

	for _, review := range reviews {
		output = append(output, ReviewOutputDTO{
			ID: review.ID,
			Author: dto.UserDTO{
				ID:       review.Author.ID,
				Username: review.Author.Username,
			},
			Rate:        review.Rate,
			Description: review.Description,
			Match:       fmt.Sprintf("http://localhost:80/football/%d", review.MatchID),
			CreatedAt:   review.CreatedAt,
		})
	}

	return output, nil
}
