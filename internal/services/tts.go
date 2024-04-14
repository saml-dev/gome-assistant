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

	tts.conn.WriteMessage(req)
}

// Say something using text-to-speech on a media player with cloud.
// Takes an entityId and an optional
// map that is translated into service_data.
func (tts TTS) CloudSay(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(tts.conn, entityId)
	req.Domain = "tts"
	req.Service = "cloud_say"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	tts.conn.WriteMessage(req)
}

// Say something using text-to-speech on a media player with google_translate.
// Takes an entityId and an optional
// map that is translated into service_data.
func (tts TTS) GoogleTranslateSay(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(tts.conn, entityId)
	req.Domain = "tts"
	req.Service = "google_translate_say"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	tts.conn.WriteMessage(req)
}
