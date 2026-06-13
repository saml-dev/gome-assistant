package services

import "context"

/* Structs */

type MediaPlayer struct {
	api API
}

/* Public API */

// Send the media player the command to clear players playlist.
// Takes entity IDs.
func (mp MediaPlayer) ClearPlaylist(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "clear_playlist",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Group players together. Only works on platforms with support for
// player groups. Takes entity IDs and an optional service_data,
// which must be serializable to a JSON object.
func (mp MediaPlayer) Join(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "join",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for next track.
// Takes entity IDs.
func (mp MediaPlayer) Next(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_next_track",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for pause.
// Takes entity IDs.
func (mp MediaPlayer) Pause(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_pause",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for play.
// Takes entity IDs.
func (mp MediaPlayer) Play(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_play",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle media player play/pause state.
// Takes entity IDs.
func (mp MediaPlayer) PlayPause(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_play_pause",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for previous track.
// Takes entity IDs.
func (mp MediaPlayer) Previous(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_previous_track",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command to seek in current playing media.
// Takes entity IDs and an optional service_data, which must be
// serializable to a JSON object.
func (mp MediaPlayer) Seek(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "media_seek",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the stop command.
// Takes entity IDs.
func (mp MediaPlayer) Stop(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_stop",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for playing media. Takes an
// entity IDs and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) PlayMedia(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "play_media",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set repeat mode. Takes entity IDs and an optional service_data,
// which must be serializable to a JSON object.
func (mp MediaPlayer) RepeatSet(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "repeat_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command to change sound mode. Takes an
// entity IDs and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) SelectSoundMode(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "select_sound_mode",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command to change input source. Takes an
// entity IDs and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) SelectSource(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "select_source",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set shuffling state. Takes entity IDs and an optional
// service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) Shuffle(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "shuffle_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggles a media player power state.
// Takes entity IDs.
func (mp MediaPlayer) Toggle(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "toggle",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player power off.
// Takes entity IDs.
func (mp MediaPlayer) TurnOff(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "turn_off",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player power on.
// Takes entity IDs.
func (mp MediaPlayer) TurnOn(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "turn_on",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
// Takes entity IDs.
func (mp MediaPlayer) Unjoin(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "unjoin",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player volume down.
// Takes entity IDs.
func (mp MediaPlayer) VolumeDown(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "volume_down",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Mute a media player's volume. Takes entity IDs and an optional
// service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) VolumeMute(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "volume_mute",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set a media player's volume level. Takes entity IDs and an
// optional service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) VolumeSet(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "volume_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player volume up.
// Takes entity IDs.
func (mp MediaPlayer) VolumeUp(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "volume_up",
		Target:  target,
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
