package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
		ctx, "tts", "clear_cache", nil, ga.Target{},
	)
}

// Say something using text-to-speech on a media player with cloud.
func (tts TTS) CloudSay(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return tts.service.CallService(
		ctx, "tts", "cloud_say",
		serviceData, target,
	)
}

// Say something using text-to-speech on a media player with
// google_translate.
func (tts TTS) GoogleTranslateSay(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return tts.service.CallService(
		ctx, "tts", "google_translate_say",
		serviceData, target,
	)
}
