package domain

type Profile struct {
	UserID         string
	Username       string
	FavoriteTeam   int64
	FollowersCount int
	FollowingCount int
	IsFollowing    bool
}

func NewProfile(userID, username string, favoriteTeam int64, followersCount, followingCount int, isFollowing bool) *Profile {
	return &Profile{
		UserID:         userID,
		Username:       username,
		FavoriteTeam:   favoriteTeam,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		IsFollowing:    isFollowing,
	}
}
