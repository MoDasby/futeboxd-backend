package football

type Match struct {
	ID    int64 `json:"id"`
	Venue struct {
		FullName string `json:"fullName"`
		City     string `json:"city"`
	} `json:"venue"`
	Date        string `json:"date"`
	Note        string `json:"note"`
	Competitors []*struct {
		HomeAway string `json:"homeAway"`
		Winner   bool   `json:"winner"`
		Score    int32  `json:"score"`
		Team     *struct {
			ID           int64  `json:"id"`
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
			Color        string `json:"color"`
			Logo         string `json:"logo"`
		} `json:"team"`
	} `json:"competitors"`
	Completed       bool   `json:"completed"`
	StatusName      string `json:"status_name"`
	CompetitionName string `json:"competition_name"`
}

type Team struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
}
