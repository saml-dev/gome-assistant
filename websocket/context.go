package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
)

type Context struct {
	ID       *string `json:"id"`
	UserID   *string `json:"user_id"`
	ParentID *string `json:"parent_id"`
}

func (c *Context) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	if b[0] == '"' {
		// The context is stored as a naked string. I think this can
		// happen but I don't know what it's supposed to signify.
		slog.Info("bare string as context; ignored", "input", string(b))
		return nil
	}

	// Unmarshal into a type that is assignable to Context but without
	// an `UnmarshalJSON()` method:
	var context struct {
		ID       *string `json:"id"`
		UserID   *string `json:"user_id"`
		ParentID *string `json:"parent_id"`
	}
	if err := json.Unmarshal(b, &context); err != nil {
		return fmt.Errorf("unmarshaling context '%s': %w", string(b), err)
	}
	*c = context
	return nil
}
