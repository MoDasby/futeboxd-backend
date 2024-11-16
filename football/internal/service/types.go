package service

type EspnTeam struct {
	ID           string `json:"id"`
	Name         string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
	Color        string `json:"color"`
	Logos        []struct {
		Href   string   `json:"href"`
		Alt    string   `json:"alt"`
		Rel    []string `json:"rel"`
		Width  int      `json:"width"`
		Height int      `json:"height"`
	} `json:"logos"`
}

type EspnTeams struct {
	Sports []struct {
		Leagues []struct {
			Teams []struct {
				Team EspnTeam
			} `json:"teams"`
		} `json:"leagues"`
	} `json:"sports"`
}

type EspnLeague struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Logos        []struct {
		Href   string   `json:"href"`
		Alt    string   `json:"alt"`
		Rel    []string `json:"rel"`
		Width  int      `json:"width"`
		Height int      `json:"height"`
	} `json:"logos"`
	HasStandings bool   `json:"hasStandings"`
	Slug         string `json:"slug"`
}

type EspnLeagues struct {
	Leagues []EspnLeague `json:"leagues"`
}

type EspnSchedule struct {
	Season struct {
		DisplayName string `json:"displayName"`
	} `json:"requestedSeason"`
	Events []struct {
		Id         string `json:"id"`
		Date       string `json:"date"`
		SeasonType struct {
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
		} `json:"seasonType"`
		Competitions []struct {
			Venue struct {
				FullName string `json:"fullName"`
				Address  struct {
					City string `json:"city"`
				}
			} `json:"venue"`
			Competitors []struct {
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Score    struct {
					Value float64 `json:"value"`
				} `json:"score"`
				Team struct {
					Id           string `json:"id"`
					Name         string `json:"displayName"`
					Abbreviation string `json:"abbreviation"`
					Color        string `json:"color"`
					Logos        []struct {
						Href   string   `json:"href"`
						Alt    string   `json:"alt"`
						Rel    []string `json:"rel"`
						Width  int      `json:"width"`
						Height int      `json:"height"`
					} `json:"logos"`
					StandingSummary string `json:"standingSummary"`
				} `json:"team"`
			} `json:"competitors"`
			Notes []struct {
				Headline string `json:"headline"`
			} `json:"notes"`
			Status struct {
				Type struct {
					Name      string `json:"name"`
					Completed bool   `json:"completed"`
				} `json:"type"`
			} `json:"status"`
		} `json:"competitions"`
	} `json:"events"`
}

type EspnEventSummary struct {
	GameInfo struct {
		Venue struct {
			FullName string `json:"fullName"`
			Address  struct {
				City    string `json:"city"`
				Country string `json:"country"`
			} `json:"address"`
		} `json:"venue"`
		Attendance int `json:"attendance"`
	} `json:"gameInfo"`
	Header struct {
		ID     string `json:"id"`
		Season struct {
			Name string `json:"name"`
		} `json:"season"`
		Competitions []struct {
			Date        string `json:"date"`
			Competitors []struct {
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Score    string `json:"score"`
				Team     struct {
					Id           string `json:"id"`
					Name         string `json:"displayName"`
					Abbreviation string `json:"abbreviation"`
					Color        string `json:"color"`
					Logos        []struct {
						Href   string   `json:"href"`
						Alt    string   `json:"alt"`
						Rel    []string `json:"rel"`
						Width  int      `json:"width"`
						Height int      `json:"height"`
					} `json:"logos"`
				} `json:"team"`
			} `json:"competitors"`
			Notes []struct {
				Headline string `json:"headline"`
			} `json:"notes"`
			Status struct {
				Type struct {
					Name      string `json:"name"`
					Completed bool   `json:"completed"`
				} `json:"type"`
			} `json:"status"`
		} `json:"competitions"`
	} `json:"header"`
	Rosters []struct {
		Roster []struct {
			Starter   bool   `json:"starter"`
			Jersey    string `json:"jersey"`
			SubbedIn  bool   `json:"subbedIn"`
			SubbedOut bool   `json:"subbedOut"`
			Athlete   struct {
				ID          string `json:"id"`
				LastName    string `json:"lastName"`
				FullName    string `json:"fullName"`
				DisplayName string `json:"displayName"`
				HeadShot    struct {
					Href string `json:"href"`
					Alt  string `json:"alt"`
				} `json:"headshot"`
			} `json:"athlete"`
			Position struct {
				DisplayName  string `json:"displayName"`
				Abbreviation string `json:"abbreviation"`
			} `json:"position"`
		} `json:"roster"`
	} `json:"rosters"`
	KeyEvents []struct {
		Type struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"type"`
		Text  string `json:"text"`
		Clock struct {
			Value        float64 `json:"value"`
			DisplayValue string  `json:"displayValue"`
		} `json:"clock"`
	} `json:"keyEvents"`
}
