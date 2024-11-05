package football

type Event struct {
	Type struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"type"`
	Text  string `json:"text"`
	Clock struct {
		Value        float64 `json:"value"`
		DisplayValue string  `json:"displayValue"`
	} `json:"clock"`
}

type Venue struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type Roster struct {
	Starter bool `json:"starter"`
	Athlete struct {
		ID          string `json:"id"`
		LastName    string `json:"lastName"`
		FullName    string `json:"fullName"`
		DisplayName string `json:"displayName"`
		HeadShot    struct {
			Href string `json:"href"`
			Alt  string `json:"alt"`
		} `json:"headshot,omitempty"`
	} `json:"athlete"`
	Position struct {
		DisplayName  string `json:"displayName"`
		Abbreviation string `json:"abbreviation"`
		SubbedIn     bool   `json:"subbedIn"`
		SubbedOut    bool   `json:"subbedOut"`
	} `json:"position"`
}

type Competitor struct {
	HomeAway string   `json:"homeAway"`
	Winner   bool     `json:"winner"`
	Score    int32    `json:"score"`
	Team     Team     `json:"team"`
	Roster   []Roster `json:"roster,omitempty"`
}

type Match struct {
	ID              int64      `json:"id"`
	Venue           Venue      `json:"venue"`
	Date            string     `json:"date"`
	Note            string     `json:"note"`
	Completed       bool       `json:"completed"`
	StatusName      string     `json:"status_name"`
	CompetitionName string     `json:"competition_name"`
	Events          []Event    `json:"events,omitempty"`
	HomeCompetitor  Competitor `json:"home_competitor"`
	AwayCompetitor  Competitor `json:"away_competitor"`
}

type Team struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
}
