package usecase

import (
	"fmt"

	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/pkg/client/user"
	"github.com/modasby/futeboxd-api/pkg/utils"
	"github.com/modasby/futeboxd-api/services/reviews/internal/domain"
)

type ListReviewsUseCase struct {
	reviewRepository domain.ReviewRepository
	footballService  football.Client
	userClient       user.Client
}

func NewListReviewsUseCase(
	reviewRepository domain.ReviewRepository,
	footballService football.Client,
	userClient user.Client,
) *ListReviewsUseCase {
	return &ListReviewsUseCase{
		reviewRepository: reviewRepository,
		footballService:  footballService,
		userClient:       userClient,
	}
}

type ReviewOutputDTO struct {
	AuthorUsername string          `json:"author_username"`
	Rate           int             `json:"rate"`
	Description    string          `json:"description"`
	Match          *football.Match `json:"match"`
}

func (uc ListReviewsUseCase) Execute(size, page int, username, team, match string) (*utils.Pageable[ReviewOutputDTO], error) {
	var userID string

	if username != "" {
		user, err := uc.userClient.GetUser(username)
		if err != nil {
			return nil, err
		}

		userID = user.ID
	}

	reviews, err := uc.reviewRepository.ListReviews(size, page, userID, team, match)
	if err != nil {
		return nil, err
	}

	if len(reviews.Items) == 0 {
		return &utils.Pageable[ReviewOutputDTO]{
			Items:       make([]ReviewOutputDTO, len(reviews.Items)),
			Size:        reviews.Size,
			TotalItems:  reviews.TotalItems,
			CurrentPage: reviews.CurrentPage,
			Pages:       reviews.Pages,
		}, nil
	}

	usersIds := func(reviews []*domain.Review) []string {
		var res []string

		for _, review := range reviews {
			res = append(res, review.UserID)
		}

		return res
	}(reviews.Items)

	users, err := uc.userClient.GetUsers(usersIds)
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]user.UserOutputDTO)

	for _, u := range users {
		usersMap[u.ID] = u
	}

	output := &utils.Pageable[ReviewOutputDTO]{
		Items:       make([]ReviewOutputDTO, len(reviews.Items)),
		Size:        reviews.Size,
		TotalItems:  reviews.TotalItems,
		CurrentPage: reviews.CurrentPage,
		Pages:       reviews.Pages,
	}

	for i := 0; i < len(reviews.Items); i++ {

		match, err := uc.footballService.GetMatch(reviews.Items[i].MatchID)
		fmt.Println(match)
		if err != nil {
			return nil, err
		}
		output.Items[i] = ReviewOutputDTO{
			AuthorUsername: usersMap[reviews.Items[i].UserID].Username, // Assumindo que user não é nil
			Rate:           reviews.Items[i].Rate,
			Description:    reviews.Items[i].Description,
			Match:          match, // Assumindo que match não é nil
		}
	}

	output.Size = reviews.Size
	output.TotalItems = reviews.TotalItems
	output.CurrentPage = reviews.CurrentPage
	output.Pages = reviews.Pages

	return output, nil
}
