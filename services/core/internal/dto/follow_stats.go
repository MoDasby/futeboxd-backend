package dto

type FollowStatsDTO struct {
	FollowersCount int  `json:"followers_count"`
	FollowingCount int  `json:"following_count"`
	Following      bool `json:"following"`
}
