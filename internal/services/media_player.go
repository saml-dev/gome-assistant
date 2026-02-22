package services

/* Structs */

type MediaPlayer struct {
	api API
}

/* Public API */

// Send the media player the command to clear players playlist.
// Takes an entityID.
func (mp MediaPlayer) ClearPlaylist(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "clear_playlist",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Group players together. Only works on platforms with support for
// player groups. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (mp MediaPlayer) Join(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "join",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Send the media player the command for next track.
// Takes an entityID.
func (mp MediaPlayer) Next(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_next_track",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Send the media player the command for pause.
// Takes an entityID.
func (mp MediaPlayer) Pause(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_pause",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Send the media player the command for play.
// Takes an entityID.
func (mp MediaPlayer) Play(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_play",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Toggle media player play/pause state.
// Takes an entityID.
func (mp MediaPlayer) PlayPause(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_play_pause",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Send the media player the command for previous track.
// Takes an entityID.
func (mp MediaPlayer) Previous(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_previous_track",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Send the media player the command to seek in current playing media.
// Takes an entityID and an optional service_data, which must be
// serializable to a JSON object.
func (mp MediaPlayer) Seek(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "media_seek",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Send the media player the stop command.
// Takes an entityID.
func (mp MediaPlayer) Stop(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "media_stop",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Send the media player the command for playing media. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) PlayMedia(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "play_media",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Set repeat mode. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (mp MediaPlayer) RepeatSet(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "repeat_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Send the media player the command to change sound mode. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) SelectSoundMode(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "select_sound_mode",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Send the media player the command to change input source. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (mp MediaPlayer) SelectSource(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "select_source",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Set shuffling state. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) Shuffle(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "shuffle_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Toggles a media player power state.
// Takes an entityID.
func (mp MediaPlayer) Toggle(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Turn a media player power off.
// Takes an entityID.
func (mp MediaPlayer) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Turn a media player power on.
// Takes an entityID.
func (mp MediaPlayer) TurnOn(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
// Takes an entityID.
func (mp MediaPlayer) Unjoin(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "unjoin",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Turn a media player volume down.
// Takes an entityID.
func (mp MediaPlayer) VolumeDown(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "volume_down",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}

// Mute a media player's volume. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) VolumeMute(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "volume_mute",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Set a media player's volume level. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (mp MediaPlayer) VolumeSet(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "media_player",
		Service:     "volume_set",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return mp.api.CallAndForget(req)
}

// Turn a media player volume up.
// Takes an entityID.
func (mp MediaPlayer) VolumeUp(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "media_player",
		Service: "volume_up",
		Target:  Entity(entityID),
	}
	return mp.api.CallAndForget(req)
}
