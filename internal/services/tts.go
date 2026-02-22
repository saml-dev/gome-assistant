package services

/* Structs */

type TTS struct {
	api API
}

/* Public API */

// Remove all text-to-speech cache files and RAM cache.
func (tts TTS) ClearCache() error {
	req := BaseServiceRequest{
		Domain:  "tts",
		Service: "clear_cache",
		Target:  Entity(""),
	}
	return tts.api.CallAndForget(req)
}

// Say something using text-to-speech on a media player with cloud.
// Takes an entityID and an optional service_data, which must be
// serializable to a JSON object.
func (tts TTS) CloudSay(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "tts",
		Service:     "cloud_say",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return tts.api.CallAndForget(req)
}

// Say something using text-to-speech on a media player with
// google_translate. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (tts TTS) GoogleTranslateSay(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "tts",
		Service:     "google_translate_say",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return tts.api.CallAndForget(req)
}
