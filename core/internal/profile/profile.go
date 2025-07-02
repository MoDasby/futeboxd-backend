package profile

type Profile struct {
	UserID           string
	Username         string
	Name             string
	Bio              string
	ProfilePicture   string
	FavoriteTeam     int64
	FollowersCount   int
	FollowingCount   int
	ReviewsCount     int
	IsFollowing      bool
	MostReviewedTeam int64
	AvgRating        float64
}
