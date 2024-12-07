package null

import "encoding/json"

type String string

func NewString(str string) String {
	return String(str)
}

func (s String) MarshalJSON() ([]byte, error) {
	if s == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(s))
}
