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
	}

	var result any
	if err := tts.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Say something using text-to-speech on a media player with cloud.
// Takes entity IDs and an optional service_data, which must be
// serializable to a JSON object.
func (tts TTS) CloudSay(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "tts",
		Service:     "cloud_say",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := tts.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Say something using text-to-speech on a media player with
// google_translate. Takes entity IDs and an optional service_data,
// which must be serializable to a JSON object.
func (tts TTS) GoogleTranslateSay(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "tts",
		Service:     "google_translate_say",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := tts.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
