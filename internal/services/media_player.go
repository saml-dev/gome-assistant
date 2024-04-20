package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type MediaPlayer struct {
	conn *websocket.Conn
}

func NewMediaPlayer(conn *websocket.Conn) *MediaPlayer {
	return &MediaPlayer{
		conn: conn,
	}
}

/* Public API */

// Send the media player the command to clear players playlist.
// Takes an entityID.
func (mp MediaPlayer) ClearPlaylist(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "clear_playlist",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Group players together. Only works on platforms with support for player groups.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Join(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "join",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command for next track.
// Takes an entityID.
func (mp MediaPlayer) Next(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_next_track",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command for pause.
// Takes an entityID.
func (mp MediaPlayer) Pause(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_pause",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command for play.
// Takes an entityID.
func (mp MediaPlayer) Play(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_play",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle media player play/pause state.
// Takes an entityID.
func (mp MediaPlayer) PlayPause(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_play_pause",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command for previous track.
// Takes an entityID.
func (mp MediaPlayer) Previous(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_previous_track",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command to seek in current playing media.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Seek(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_seek",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}
	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the stop command.
// Takes an entityID.
func (mp MediaPlayer) Stop(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "media_stop",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command for playing media.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) PlayMedia(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "play_media",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Set repeat mode. Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) RepeatSet(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "repeat_set",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command to change sound mode.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) SelectSoundMode(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "select_sound_mode",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the media player the command to change input source.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) SelectSource(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "select_source",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Set shuffling state.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) Shuffle(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "shuffle_set",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggles a media player power state.
// Takes an entityID.
func (mp MediaPlayer) Toggle(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Turn a media player power off.
// Takes an entityID.
func (mp MediaPlayer) TurnOff(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "turn_off",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Turn a media player power on.
// Takes an entityID.
func (mp MediaPlayer) TurnOn(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Unjoin the player from a group. Only works on
// platforms with support for player groups.
// Takes an entityID.
func (mp MediaPlayer) Unjoin(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "unjoin",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Turn a media player volume down.
// Takes an entityID.
func (mp MediaPlayer) VolumeDown(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "volume_down",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Mute a media player's volume.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) VolumeMute(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "volume_mute",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Set a media player's volume level.
// Takes an entityID and an optional
// map that is translated into service_data.
func (mp MediaPlayer) VolumeSet(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "volume_set",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Turn a media player volume up.
// Takes an entityID.
func (mp MediaPlayer) VolumeUp(entityID string) {
	req := CallServiceRequest{
		Domain:  "media_player",
		Service: "volume_up",
		Target: Target{
			EntityID: entityID,
		},
	}

	mp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
