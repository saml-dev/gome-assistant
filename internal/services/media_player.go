package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type MediaPlayer struct {
	service Service
}

func NewMediaPlayer(service Service) *MediaPlayer {
	return &MediaPlayer{
		service: service,
	}
}

/* Public API */

// Send the media player the command to clear players playlist.
func (mp MediaPlayer) ClearPlaylist(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "clear_playlist",
		nil, target,
	)
}

// Group players together. Only works on platforms with support for player groups.
func (mp MediaPlayer) Join(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "join",
		serviceData, target,
	)
}

// Send the media player the command for next track.
func (mp MediaPlayer) Next(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_next_track",
		nil, target,
	)
}

// Send the media player the command for pause.
func (mp MediaPlayer) Pause(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_pause",
		nil, target,
	)
}

// Send the media player the command for play.
func (mp MediaPlayer) Play(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_play",
		nil, target,
	)
}

// Toggle media player play/pause state.
func (mp MediaPlayer) PlayPause(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_play_pause",
		nil, target,
	)
}

// Send the media player the command for previous track.
func (mp MediaPlayer) Previous(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_previous_track",
		nil, target,
	)
}

// Send the media player the command to seek in current playing media.
func (mp MediaPlayer) Seek(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_seek",
		serviceData, target,
	)
}

// Send the media player the stop command.
func (mp MediaPlayer) Stop(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_stop",
		nil, target,
	)
}

// Send the media player the command for playing media.
func (mp MediaPlayer) PlayMedia(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "play_media",
		serviceData, target,
	)
}

// Set repeat mode.
func (mp MediaPlayer) RepeatSet(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "repeat_set",
		serviceData, target,
	)
}

// Send the media player the command to change sound mode.
func (mp MediaPlayer) SelectSoundMode(
	target ga.Target, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "select_sound_mode",
		serviceData, target,
	)
}

// Send the media player the command to change input source.
func (mp MediaPlayer) SelectSource(
	target ga.Target, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "select_source",
		serviceData, target,
	)
}

// Set shuffling state.
func (mp MediaPlayer) Shuffle(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "shuffle_set",
		serviceData, target,
	)
}

// Toggles a media player power state.
func (mp MediaPlayer) Toggle(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "toggle",
		nil, target,
	)
}

// Turn a media player power off.
func (mp MediaPlayer) TurnOff(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "turn_off",
		nil, target,
	)
}

// Turn a media player power on.
func (mp MediaPlayer) TurnOn(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "turn_on",
		nil, target,
	)
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
func (mp MediaPlayer) Unjoin(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "unjoin",
		nil, target,
	)
}

// Turn a media player volume down.
func (mp MediaPlayer) VolumeDown(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_down",
		nil, target,
	)
}

// Mute a media player's volume.
func (mp MediaPlayer) VolumeMute(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_mute",
		serviceData, target,
	)
}

// Set a media player's volume level.
func (mp MediaPlayer) VolumeSet(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_set",
		serviceData, target,
	)
}

// Turn a media player volume up.
func (mp MediaPlayer) VolumeUp(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_up",
		nil, target,
	)
}
