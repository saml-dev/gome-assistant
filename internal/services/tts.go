package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type TTS struct {
	conn *websocket.Conn
}

func NewTTS(conn *websocket.Conn) *TTS {
	return &TTS{
		conn: conn,
	}
}

/* Public API */

// Remove all text-to-speech cache files and RAM cache.
func (tts TTS) ClearCache() {
	req := NewBaseServiceRequest(tts.conn, "")
	req.Domain = "tts"
	req.Service = "clear_cache"

	tts.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Say something using text-to-speech on a media player with cloud.
// Takes an entityId and an optional
// map that is translated into service_data.
func (tts TTS) CloudSay(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(tts.conn, entityId)
	req.Domain = "tts"
	req.Service = "cloud_say"
	req.ServiceData = serviceData

	tts.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Say something using text-to-speech on a media player with google_translate.
// Takes an entityId and an optional
// map that is translated into service_data.
func (tts TTS) GoogleTranslateSay(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(tts.conn, entityId)
	req.Domain = "tts"
	req.Service = "google_translate_say"
	req.ServiceData = serviceData

	tts.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
