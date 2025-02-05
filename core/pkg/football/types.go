package football

type Team struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
}

type Competitor struct {
	HomeAway string `json:"homeAway"`
	Winner   bool   `json:"winner"`
	Score    int32  `json:"score"`
	Team     Team   `json:"team"`
}

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type Match struct {
	ID             int64      `json:"id"`
	Venue          string     `json:"venue"`
	Date           string     `json:"date"`
	Note           string     `json:"note"`
	Completed      bool       `json:"completed"`
	StatusName     string     `json:"status_name"`
	League         League     `json:"league"`
	HomeCompetitor Competitor `json:"home_competitor"`
	AwayCompetitor Competitor `json:"away_competitor"`
}
