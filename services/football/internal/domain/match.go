package domain

type Venue struct {
	FullName string `json:"fullName"`
	City     string `json:"city"`
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
	HomeAway string    `json:"homeAway"`
	Winner   bool      `json:"winner"`
	Score    int32     `json:"score"`
	Team     *Team     `json:"team"`
	Roster   []*Roster `json:"roster,omitempty"`
}

type Match struct {
	ID              int64         `json:"id"`
	Venue           Venue         `json:"venue"`
	Date            string        `json:"date"`
	Note            string        `json:"note"`
	Competitors     []*Competitor `json:"competitors"`
	Completed       bool          `json:"completed"`
	StatusName      string        `json:"status_name"`
	CompetitionName string        `json:"competition_name"`
}

func NewMatch(
	ID int64,
	venue Venue,
	attendance int,
	date string,
	note string,
	competitors []*Competitor,
	completed bool,
	statusName string,
	competitionName string,
) *Match {
	return &Match{
		ID:              ID,
		Venue:           venue,
		Date:            date,
		Note:            note,
		Competitors:     competitors,
		Completed:       completed,
		StatusName:      statusName,
		CompetitionName: competitionName,
	}
}

func NewRoster(starter bool, athleteID, athleteLastName, athleteFullName, athleteDisplayName, athleteHeadShotHref, athleteHeadShotAlt, positionDisplayName, positionAbbreviation string, subbedIn, subbedOut bool) *Roster {
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

func NewVenue(fullName string, city string) *Venue {
	return &Venue{
		FullName: fullName, City: city,
	}
}

func NewCompetitor(homeAway string, winner bool, score int32, team *Team, roster []*Roster) *Competitor {
	return &Competitor{
		HomeAway: homeAway, Winner: winner, Score: score, Team: team, Roster: roster,
	}
}
