package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type ListFeedUsecase struct {
	reviewsRepo    domain.ReviewRepository
	userRepository domain.UserRepository
	footballClient football.Client
}

func NewListFeedUsecase(
	reviewsRepo domain.ReviewRepository,
	userRepository domain.UserRepository,
	footballClient football.Client,
) *ListFeedUsecase {
	return &ListFeedUsecase{
		reviewsRepo:    reviewsRepo,
		userRepository: userRepository,
		footballClient: footballClient,
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

	matches, err := GetMatchesByReview(uc.footballClient, reviews)
	if err != nil {
		return nil, err
	}

	output := make([]dto.ReviewDTO, 0)

	for _, review := range reviews {
		match := matches[review.MatchID]

		output = append(output, dto.ReviewDTO{
			ID: review.ID,
			Author: dto.ContentAuthorDTO{
				ID:       review.Author.ID,
				Username: review.Author.Username,
			},
			Rate:          review.Rate,
			Description:   review.Description,
			Match:         &match,
			Likes:         review.Likes,
			CommentsCount: review.CommentsCount,
			IsLiked:       review.IsLiked,
			CreatedAt:     review.CreatedAt,
		})
	}

	return output, nil
}
