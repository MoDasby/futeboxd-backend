package usecase

import (
	"fmt"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type ListReviewsUseCase struct {
	reviewRepository domain.ReviewRepository
	userRepository   domain.UserRepository
}

func NewListReviewsUseCase(
	reviewRepository domain.ReviewRepository,
	userRepository domain.UserRepository,
) *ListReviewsUseCase {
	return &ListReviewsUseCase{
		reviewRepository: reviewRepository,
		userRepository:   userRepository,
	}
}

func (uc ListReviewsUseCase) Execute(
	requesterID,
	username,
	team,
	match string,
	page *pagination.Page,
) ([]dto.ReviewDTO, error) {

	userID := ""

	if username != "" {
		user, err := uc.userRepository.FindOneByIdOrUsername(username)
		if err != nil {
			return nil, err
		}

		userID = user.ID
	}

	reviews, err := uc.reviewRepository.ListAll(requesterID, userID, team, match, page)
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
