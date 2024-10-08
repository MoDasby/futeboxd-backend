package domain

import "github.com/modasby/futeboxd-api/pkg/utils"

type ReviewRepository interface {
	AddReview(review *Review) error
	ListReviews(pageSize, pageIndex int, userID, team, match string) (*utils.Pageable[*Review], error)
}
