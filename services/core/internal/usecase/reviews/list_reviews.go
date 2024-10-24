package usecase

import (
	"fmt"
	"time"

	"github.com/modasby/futeboxd-api/pkg/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
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

type ReviewOutputDTO struct {
	ID          int         `json:"id"`
	Author      dto.UserDTO `json:"author"`
	Rate        int         `json:"rate"`
	Description string      `json:"description"`
	Match       string      `json:"match"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (uc ListReviewsUseCase) Execute(
	username,
	team,
	match string,
	page pagination.Page,
) ([]ReviewOutputDTO, error) {

	userID := ""

	if username != "" {
		user, err := uc.userRepository.FindOneByIdOrUsername(username)
		if err != nil {
			return nil, err
		}

		userID = user.ID
	}

	reviews, err := uc.reviewRepository.ListAll(page.Size, page.Index, userID, team, match)
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
