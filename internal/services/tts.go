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
	return tts.api.Call(req)
}

// Say something using text-to-speech on a media player with cloud.
// Takes an entityId and an optional
// map that is translated into service_data.
func (tts TTS) CloudSay(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "tts",
		Service: "cloud_say",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return tts.api.Call(req)
}

// Say something using text-to-speech on a media player with google_translate.
// Takes an entityId and an optional
// map that is translated into service_data.
func (tts TTS) GoogleTranslateSay(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "tts",
		Service: "google_translate_say",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return tts.api.Call(req)
}
