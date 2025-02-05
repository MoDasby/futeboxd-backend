package models

import (
	"encoding/json"
)

type Events []Event

func (e Events) MarshalJSON() ([]byte, error) {
	if len(e) == 0 {
		return []byte("[]"), nil // Força um array vazio
	}
	return json.Marshal([]Event(e)) // Serializa normalmente
}

type MatchSummary struct {
	Events Events `json:"events"`
}
