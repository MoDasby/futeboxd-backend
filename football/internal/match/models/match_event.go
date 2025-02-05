package models

type Event struct {
	Type struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"type"`
	ClockValue      int    `json:"clock_value"`
	TeamID          int    `json:"team_id,omitempty"`
	ParticipantName string `json:"participant_name,omitempty"`
}

func NewEvent(
	typeID,
	typeText string,
	clockValue int,
	teamID int,
	participantName string,
) *Event {
	return &Event{
		Type: struct {
			ID   string "json:\"id\""
			Text string "json:\"text\""
		}{
			ID: typeID, Text: typeText,
		},
		ClockValue:      clockValue,
		TeamID:          teamID,
		ParticipantName: participantName,
	}
}
