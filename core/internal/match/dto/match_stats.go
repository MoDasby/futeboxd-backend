package dto

import "github.com/modasby/futeboxd-backend/core/pkg/football"

type ReviewsSummary struct {
	Count1 float32 `json:"count1"`
	Count2 float32 `json:"count2"`
	Count3 float32 `json:"count3"`
	Count4 float32 `json:"count4"`
	Count5 float32 `json:"count5"`
}

type MatchStats struct {
	Match          *football.Match `json:"match"`
	AvgRate        float32         `json:"average_rate"`
	ReviewsSummary ReviewsSummary  `json:"reviews_summary"`
}
