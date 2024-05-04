package websocket

import (
	"encoding/json"
	"fmt"
)

// RawMessage is like `json.RawMessage`, except that its `String()`
// method converts it directly to a string.
type RawMessage json.RawMessage

// UnmarshalJSON delegates to `json.RawMessage`. (The method has a
// pointer receiver, so we have to implement it explicitly.)
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(m).UnmarshalJSON(data)
}

func (rm RawMessage) String() string {
	return string(rm)
}

// BaseMessage implements the required part of any websocket message.
// The idea is to embed this type in other message types.
type BaseMessage struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

type Request interface {
	GetID() int64
	SetID(id int64)
}

func (msg *BaseMessage) GetID() int64 {
	return msg.ID
}

func (msg *BaseMessage) SetID(id int64) {
	msg.ID = id
}

// Message holds a complete message, only partly parsed. The entire,
// original, unparsed message is available in the `Raw` field.
type Message struct {
	BaseMessage

	// Raw contains the original, full, unparsed message (including
	// fields `Type` and `ID`, which also appear in `BaseMessage`).
	Raw RawMessage `json:"-"`
}

type BaseResultMessage struct {
	BaseMessage
	Success bool         `json:"success"`
	Error   *ResultError `json:"error,omitempty"`
}

type ResultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *ResultError) Error() string {
	switch {
	case err.Code != "" && err.Message != "":
		return fmt.Sprintf("%s: %s", err.Code, err.Message)
	case err.Code == "" && err.Message != "":
		return fmt.Sprintf("unknown_error: %s", err.Message)
	case err.Code != "" && err.Message == "":
		return fmt.Sprintf("%s", err.Code)
	default:
		// This seems not to be an error at all.
		return fmt.Sprintf("INVALID (seems not to be an error)")
	}
}

type ResultMessage struct {
	BaseResultMessage

	// Raw contains the "result" part of the message, unparsed.
	Result RawMessage `json:"result"`
}

// GetResult parses a result out of `msg` (which must have type
// "result"). If `msg` indicates that an error occurred, convert that
// to an error and return it. Parse the result into `result`, which
// must be unmarshalable as JSON.
func (msg Message) GetResult(result any) error {
	if msg.Type != "result" {
		return fmt.Errorf(
			"response message was not of type 'result': %#v", msg,
		)
	}
	var resultMsg ResultMessage
	if err := json.Unmarshal(msg.Raw, &resultMsg); err != nil {
		return fmt.Errorf("unmarshaling result message: %w", err)
	}
	if !resultMsg.Success {
		if resultMsg.Error == nil {
			return fmt.Errorf(
				"request did not succeed but no error was returned",
			)
		}
		return resultMsg.Error
	}

	if err := json.Unmarshal(resultMsg.Result, result); err != nil {
		return fmt.Errorf("unmarshalling result from %q: %w", resultMsg.Result, err)
	}
	return nil
}
