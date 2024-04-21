package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type TTS struct {
	service Service
}

func NewTTS(service Service) *TTS {
	return &TTS{
		service: service,
	}
}

/* Public API */

// Remove all text-to-speech cache files and RAM cache.
func (tts TTS) ClearCache() (websocket.Message, error) {
	ctx := context.TODO()
	return tts.service.CallService(
		ctx, "tts", "clear_cache", nil, Target{},
	)
}

// Say something using text-to-speech on a media player with cloud.
func (tts TTS) CloudSay(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return tts.service.CallService(
		ctx, "tts", "cloud_say",
		serviceData, EntityTarget(entityID),
	)
}

// Say something using text-to-speech on a media player with
// google_translate.
func (tts TTS) GoogleTranslateSay(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return tts.service.CallService(
		ctx, "tts", "google_translate_say",
		serviceData, EntityTarget(entityID),
	)
}
