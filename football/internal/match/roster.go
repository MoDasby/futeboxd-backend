package match

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
