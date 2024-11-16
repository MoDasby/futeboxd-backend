package dto

type ProfileDTO struct {
	*UserDTO
	FollowersCount int  `json:"followers_count"`
	FollowingCount int  `json:"following_count"`
	Following      bool `json:"following"`
	Self           bool `json:"self"`
}
