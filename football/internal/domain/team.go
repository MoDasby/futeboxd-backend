package domain

type Team struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Color        string `json:"color"`
	Logo         string `json:"logo"`
}

func NewTeam(ID int64, name string, abbreviation string, color string, logo string) *Team {
	return &Team{
		ID: ID, Name: name, Abbreviation: abbreviation, Color: color, Logo: logo,
	}
}
