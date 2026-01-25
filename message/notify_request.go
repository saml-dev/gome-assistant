package message

type NotifyData struct {
	Message string `json:"message"`
	Title   string `json:"title,omitzero"`
	Data    any    `json:"data,omitzero"`
}
