package profile

type Profile struct {
	UserID         string
	Username       string
	Name           string
	Bio            string
	FavoriteTeam   int64
	FollowersCount int
	FollowingCount int
	IsFollowing    bool
}
