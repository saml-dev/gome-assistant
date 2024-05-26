package websocket

import (
	"encoding/json"
)

// RawMessage is like `json.RawMessage`, except that its `String()`
// method converts it directly to a string.
type RawMessage json.RawMessage

func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON delegates to `json.RawMessage`. (The method has a
// pointer receiver, so we have to implement it explicitly.)
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(m).UnmarshalJSON(data)
}

func (rm RawMessage) String() string {
	return string(rm)
}
