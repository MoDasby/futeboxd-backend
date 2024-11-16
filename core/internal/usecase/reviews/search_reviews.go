package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type SearchReviews struct {
	reviewsRepo    domain.ReviewRepository
	footballClient football.Client
}

func NewSearchReviews(
	reviewsRepo domain.ReviewRepository,
	footballClient football.Client,
) *SearchReviews {
	return &SearchReviews{
		reviewsRepo:    reviewsRepo,
		footballClient: footballClient,
	}
}

func (uc *SearchReviews) Execute(requesterID, term string, page *pagination.Page) ([]dto.ReviewDTO, error) {
	reviews, err := uc.reviewsRepo.Search(requesterID, term, page)
	if err != nil {
		return nil, err
	}

	matches := make(map[int64]football.Match, 0)

	if len(reviews) > 0 {
		matches, err = GetMatchesByReview(uc.footballClient, reviews)
		if err != nil {
			return nil, err
		}
	}

	output := make([]dto.ReviewDTO, len(reviews))

	for index, review := range reviews {
		match := matches[review.MatchID]

		output[index] = dto.ReviewDTO{
			ID: review.ID,
			Author: dto.UserDTO{
				ID:       review.Author.ID,
				Username: review.Author.Username,
			},
			Rate:        review.Rate,
			Description: review.Description,
			Match:       &match,
			Likes:       review.Likes,
			IsLiked:     review.IsLiked,
			CreatedAt:   review.CreatedAt,
		}
	}

	return output, nil
}
