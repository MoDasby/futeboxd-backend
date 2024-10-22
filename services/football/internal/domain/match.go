package domain

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

func NewMatch(
	ID int64,
	venue Venue,
	date string,
	note string,
	homeCompetitor Competitor,
	awayCompetitor Competitor,
	completed bool,
	statusName string,
	competitionName string,
	events []Event,
) *Match {
	return &Match{
		ID:              ID,
		Venue:           venue,
		Date:            date,
		Note:            note,
		HomeCompetitor:  homeCompetitor,
		AwayCompetitor:  awayCompetitor,
		Completed:       completed,
		StatusName:      statusName,
		CompetitionName: competitionName,
		Events:          events,
	}
}

func NewRoster(
	starter bool,
	athleteID,
	athleteLastName,
	athleteFullName,
	athleteDisplayName,
	athleteHeadShotHref,
	athleteHeadShotAlt,
	positionDisplayName,
	positionAbbreviation string,
	subbedIn,
	subbedOut bool,
) *Roster {
	return &Roster{
		Starter: starter,
		Athlete: struct {
			ID          string `json:"id"`
			LastName    string `json:"lastName"`
			FullName    string `json:"fullName"`
			DisplayName string `json:"displayName"`
			HeadShot    struct {
				Href string `json:"href"`
				Alt  string `json:"alt"`
			} `json:"headshot,omitempty"`
		}{
			ID:          athleteID,
			LastName:    athleteLastName,
			FullName:    athleteFullName,
			DisplayName: athleteDisplayName,
			HeadShot: struct {
				Href string `json:"href"`
				Alt  string `json:"alt"`
			}{
				Href: athleteHeadShotHref,
				Alt:  athleteHeadShotAlt,
			},
		},
		Position: struct {
			DisplayName  string `json:"displayName"`
			Abbreviation string `json:"abbreviation"`
			SubbedIn     bool   `json:"subbedIn"`
			SubbedOut    bool   `json:"subbedOut"`
		}{
			DisplayName:  positionDisplayName,
			Abbreviation: positionAbbreviation,
			SubbedIn:     subbedIn,
			SubbedOut:    subbedOut,
		},
	}
}

func NewVenue(name string, city string) *Venue {
	return &Venue{
		Name: name, City: city,
	}
}

func NewCompetitor(homeAway string, winner bool, score int32, team Team, roster []Roster) *Competitor {
	return &Competitor{
		HomeAway: homeAway, Winner: winner, Score: score, Team: team, Roster: roster,
	}
}

func NewEvent(
	typeID,
	typeText,
	text string,
	clockValue float64,
	clockDisplayValue string,
) *Event {
	return &Event{
		Type: struct {
			ID   string "json:\"id\""
			Text string "json:\"text\""
		}{
			ID: typeID, Text: typeText,
		},
		Text: text,
		Clock: struct {
			Value        float64 "json:\"value\""
			DisplayValue string  "json:\"displayValue\""
		}{
			Value: clockValue, DisplayValue: clockDisplayValue,
		},
	}
}
