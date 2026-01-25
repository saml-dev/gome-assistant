package message

import (
	"encoding/json"
)

// RawMessage is like `json.RawMessage`, but with a `String()` method
// that returns the JSON as a string.
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

func (m RawMessage) String() string {
	return string(m)
}
