package services

/* Structs */

type MediaPlayer struct {
	api API
}

/* Public API */

// Send the media player the command to clear players playlist.
// Takes an entityId.
func (mp MediaPlayer) ClearPlaylist(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "clear_playlist"

	return mp.api.WriteMessage(req)
}

// Group players together. Only works on platforms with support for player groups.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Join(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "join"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Send the media player the command for next track.
// Takes an entityId.
func (mp MediaPlayer) Next(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_next_track"

	return mp.api.WriteMessage(req)
}

// Send the media player the command for pause.
// Takes an entityId.
func (mp MediaPlayer) Pause(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_pause"

	return mp.api.WriteMessage(req)
}

// Send the media player the command for play.
// Takes an entityId.
func (mp MediaPlayer) Play(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_play"

	return mp.api.WriteMessage(req)
}

// Toggle media player play/pause state.
// Takes an entityId.
func (mp MediaPlayer) PlayPause(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_play_pause"

	return mp.api.WriteMessage(req)
}

// Send the media player the command for previous track.
// Takes an entityId.
func (mp MediaPlayer) Previous(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_previous_track"

	return mp.api.WriteMessage(req)
}

// Send the media player the command to seek in current playing media.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Seek(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_seek"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Send the media player the stop command.
// Takes an entityId.
func (mp MediaPlayer) Stop(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "media_stop"

	return mp.api.WriteMessage(req)
}

// Send the media player the command for playing media.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) PlayMedia(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "play_media"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Set repeat mode. Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) RepeatSet(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "repeat_set"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Send the media player the command to change sound mode.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) SelectSoundMode(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "select_sound_mode"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Send the media player the command to change input source.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) SelectSource(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "select_source"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Set shuffling state.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Shuffle(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "shuffle_set"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Toggles a media player power state.
// Takes an entityId.
func (mp MediaPlayer) Toggle(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "toggle"

	return mp.api.WriteMessage(req)
}

// Turn a media player power off.
// Takes an entityId.
func (mp MediaPlayer) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "turn_off"

	return mp.api.WriteMessage(req)
}

// Turn a media player power on.
// Takes an entityId.
func (mp MediaPlayer) TurnOn(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "turn_on"

	return mp.api.WriteMessage(req)
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
// Takes an entityId.
func (mp MediaPlayer) Unjoin(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "unjoin"

	return mp.api.WriteMessage(req)
}

// Turn a media player volume down.
// Takes an entityId.
func (mp MediaPlayer) VolumeDown(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_down"

	return mp.api.WriteMessage(req)
}

// Mute a media player's volume.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) VolumeMute(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_mute"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Set a media player's volume level.
// Takes an entityId and an optional
// map that is translated into service_data.
func (mp MediaPlayer) VolumeSet(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_set"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return mp.api.WriteMessage(req)
}

// Turn a media player volume up.
// Takes an entityId.
func (mp MediaPlayer) VolumeUp(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "media_player"
	req.Service = "volume_up"

	return mp.api.WriteMessage(req)
}
