package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type MediaPlayer struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// Send the media player the command to clear players playlist.
// Takes an entityId.
func (mp MediaPlayer) ClearPlaylist(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "clear_playlist"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Group players together. Only works on platforms with support for player groups.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Join(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "join"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command for next track.
// Takes an entityId.
func (mp MediaPlayer) Next(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_next_track"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command for pause.
// Takes an entityId.
func (mp MediaPlayer) Pause(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_pause"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command for play.
// Takes an entityId.
func (mp MediaPlayer) Play(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_play"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Toggle media player play/pause state.
// Takes an entityId.
func (mp MediaPlayer) PlayPause(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_play_pause"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command for previous track.
// Takes an entityId.
func (mp MediaPlayer) Previous(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_previous_track"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command to seek in current playing media.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Seek(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_seek"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the stop command.
// Takes an entityId.
func (mp MediaPlayer) Stop(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_stop"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command for playing media.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) PlayMedia(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "play_media"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Set repeat mode. Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) RepeatSet(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "repeat_set"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command to change sound mode.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) SelectSoundMode(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "select_sound_mode"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Send the media player the command to change input source.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) SelectSource(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "select_source"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Set shuffling state.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Shuffle(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "shuffle_set"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Toggles a media player power state.
// Takes an entityId.
func (mp MediaPlayer) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "toggle"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Turn a media player power off.
// Takes an entityId.
func (mp MediaPlayer) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "turn_off"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Turn a media player power on.
// Takes an entityId.
func (mp MediaPlayer) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "turn_on"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
// Takes an entityId.
func (mp MediaPlayer) Unjoin(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "unjoin"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Turn a media player volume down.
// Takes an entityId.
func (mp MediaPlayer) VolumeDown(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_down"

	mp.conn.WriteMessage(req, mp.ctx)
}

// Mute a media player's volume.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) VolumeMute(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_mute"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Set a media player's volume level.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) VolumeSet(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_set"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	mp.conn.WriteMessage(req, mp.ctx)
}

// Turn a media player volume up.
// Takes an entityId.
func (mp MediaPlayer) VolumeUp(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_up"

	mp.conn.WriteMessage(req, mp.ctx)
}
