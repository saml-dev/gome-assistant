package services

import (
	"context"

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
func (mp MediaPlayer) ClearPlaylist(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "clear_playlist",
		nil, EntityTarget(entityID),
	)
}

// Group players together. Only works on platforms with support for player groups.
func (mp MediaPlayer) Join(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "join",
		serviceData, EntityTarget(entityID),
	)
}

// Send the media player the command for next track.
func (mp MediaPlayer) Next(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_next_track",
		nil, EntityTarget(entityID),
	)
}

// Send the media player the command for pause.
func (mp MediaPlayer) Pause(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_pause",
		nil, EntityTarget(entityID),
	)
}

// Send the media player the command for play.
func (mp MediaPlayer) Play(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_play",
		nil, EntityTarget(entityID),
	)
}

// Toggle media player play/pause state.
func (mp MediaPlayer) PlayPause(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_play_pause",
		nil, EntityTarget(entityID),
	)
}

// Send the media player the command for previous track.
func (mp MediaPlayer) Previous(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_previous_track",
		nil, EntityTarget(entityID),
	)
}

// Send the media player the command to seek in current playing media.
func (mp MediaPlayer) Seek(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_seek",
		serviceData, EntityTarget(entityID),
	)
}

// Send the media player the stop command.
func (mp MediaPlayer) Stop(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "media_stop",
		nil, EntityTarget(entityID),
	)
}

// Send the media player the command for playing media.
func (mp MediaPlayer) PlayMedia(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "play_media",
		serviceData, EntityTarget(entityID),
	)
}

// Set repeat mode.
func (mp MediaPlayer) RepeatSet(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "repeat_set",
		serviceData, EntityTarget(entityID),
	)
}

// Send the media player the command to change sound mode.
func (mp MediaPlayer) SelectSoundMode(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "select_sound_mode",
		serviceData, EntityTarget(entityID),
	)
}

// Send the media player the command to change input source.
func (mp MediaPlayer) SelectSource(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "select_source",
		serviceData, EntityTarget(entityID),
	)
}

// Set shuffling state.
func (mp MediaPlayer) Shuffle(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "shuffle_set",
		serviceData, EntityTarget(entityID),
	)
}

// Toggles a media player power state.
func (mp MediaPlayer) Toggle(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "toggle",
		nil, EntityTarget(entityID),
	)
}

// Turn a media player power off.
func (mp MediaPlayer) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "turn_off",
		nil, EntityTarget(entityID),
	)
}

// Turn a media player power on.
func (mp MediaPlayer) TurnOn(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "turn_on",
		nil, EntityTarget(entityID),
	)
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
func (mp MediaPlayer) Unjoin(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "unjoin",
		nil, EntityTarget(entityID),
	)
}

// Turn a media player volume down.
func (mp MediaPlayer) VolumeDown(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_down",
		nil, EntityTarget(entityID),
	)
}

// Mute a media player's volume.
func (mp MediaPlayer) VolumeMute(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_mute",
		serviceData, EntityTarget(entityID),
	)
}

// Set a media player's volume level.
func (mp MediaPlayer) VolumeSet(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_set",
		serviceData, EntityTarget(entityID),
	)
}

// Turn a media player volume up.
func (mp MediaPlayer) VolumeUp(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return mp.service.CallService(
		ctx, "media_player", "volume_up",
		nil, EntityTarget(entityID),
	)
}
