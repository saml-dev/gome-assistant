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
	req := CallServiceRequest{
		Domain:  "tts",
		Service: "clear_cache",
	}

	tts.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Say something using text-to-speech on a media player with cloud.
func (tts TTS) CloudSay(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "tts",
		Service: "cloud_say",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	tts.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Say something using text-to-speech on a media player with
// google_translate.
func (tts TTS) GoogleTranslateSay(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "tts",
		Service: "google_translate_say",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	tts.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
