package services

import "context"

/* Structs */

type TTS struct {
	api API
}

/* Public API */

// Remove all text-to-speech cache files and RAM cache.
func (tts TTS) ClearCache(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "tts",
		Service: "clear_cache",
		Target:  Entity(""),
	}

	var result any
	if err := tts.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Say something using text-to-speech on a media player with cloud.
// Takes an entityID and an optional service_data, which must be
// serializable to a JSON object.
func (tts TTS) CloudSay(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "tts",
		Service:     "cloud_say",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	var result any
	if err := tts.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Say something using text-to-speech on a media player with
// google_translate. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (tts TTS) GoogleTranslateSay(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "tts",
		Service:     "google_translate_say",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	var result any
	if err := tts.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
