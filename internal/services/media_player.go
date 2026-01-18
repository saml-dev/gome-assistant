package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type MediaPlayer struct {
	api API
}

/* Public API */

// Send the media player the command to clear players playlist.
// Takes an entityID.
func (mp MediaPlayer) ClearPlaylist(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "clear_playlist",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Group players together. Only works on platforms with support for
// player groups. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (mp MediaPlayer) Join(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "join",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for next track.
// Takes an entityID.
func (mp MediaPlayer) Next(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_next_track",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for pause.
// Takes an entityID.
func (mp MediaPlayer) Pause(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_pause",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for play.
// Takes an entityID.
func (mp MediaPlayer) Play(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_play",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle media player play/pause state.
// Takes an entityID.
func (mp MediaPlayer) PlayPause(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_play_pause",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for previous track.
// Takes an entityID.
func (mp MediaPlayer) Previous(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_previous_track",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command to seek in current playing media.
// Takes an entityID and an optional service_data, which must be
// serializable to a JSON object.
func (mp MediaPlayer) Seek(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "media_seek",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the stop command.
// Takes an entityID.
func (mp MediaPlayer) Stop(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_stop",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command for playing media. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) PlayMedia(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "play_media",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set repeat mode. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (mp MediaPlayer) RepeatSet(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "repeat_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command to change sound mode. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) SelectSoundMode(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "select_sound_mode",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the media player the command to change input source. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) SelectSource(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "select_source",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set shuffling state. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) Shuffle(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "shuffle_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggles a media player power state.
// Takes an entityID.
func (mp MediaPlayer) Toggle(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "toggle",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player power off.
// Takes an entityID.
func (mp MediaPlayer) TurnOff(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "turn_off",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player power on.
// Takes an entityID.
func (mp MediaPlayer) TurnOn(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "turn_on",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
// Takes an entityID.
func (mp MediaPlayer) Unjoin(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "unjoin",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player volume down.
// Takes an entityID.
func (mp MediaPlayer) VolumeDown(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "volume_down",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Mute a media player's volume. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) VolumeMute(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "volume_mute",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set a media player's volume level. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) VolumeSet(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "volume_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Turn a media player volume up.
// Takes an entityID.
func (mp MediaPlayer) VolumeUp(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "volume_up",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := mp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
