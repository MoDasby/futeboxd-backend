package match

type MatchSummary struct {
	HomeCompetitorRoster []Roster `json:"home_competitor_roster"`
	AwayCompetitorRoster []Roster `json:"away_competitor_roster"`
	Events               []Event  `json:"events"`
}
