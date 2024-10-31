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

func (uc *ListFeedUsecase) Execute(requesterID, strategy string, page *pagination.Page) ([]dto.ReviewDTO, error) {
	if strategy == "" {
		strategy = "relevant"
	}

	reviews, err := uc.reviewsRepo.ListFeed(requesterID, strategy, page)
	if err != nil {
		return nil, err
	}

	output := make([]dto.ReviewDTO, 0)

	for _, review := range reviews {
		output = append(output, dto.ReviewDTO{
			ID: review.ID,
			Author: dto.UserDTO{
				ID:       review.Author.ID,
				Username: review.Author.Username,
			},
			Rate:        review.Rate,
			Description: review.Description,
			Match:       fmt.Sprintf("http://localhost:80/football/%d", review.MatchID),
			Likes:       review.Likes,
			IsLiked:     review.IsLiked,
			CreatedAt:   review.CreatedAt,
		})
	}

	return output, nil
}
