package gomeassistant

import (
	"context"

	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/internal/services"
	ws "saml.dev/gome-assistant/internal/websocket"
)

type Service struct {
	AlarmControlPanel *services.AlarmControlPanel
	Climate           *services.Climate
	Cover             *services.Cover
	HomeAssistant     *services.HomeAssistant
	Light             *services.Light
	Lock              *services.Lock
	MediaPlayer       *services.MediaPlayer
	Switch            *services.Switch
	InputBoolean      *services.InputBoolean
	InputButton       *services.InputButton
	InputText         *services.InputText
	InputDatetime     *services.InputDatetime
	InputNumber       *services.InputNumber
	Event             *services.Event
	Notify            *services.Notify
	Number            *services.Number
	Scene             *services.Scene
	Script            *services.Script
	TTS               *services.TTS
	Vacuum            *services.Vacuum
	ZWaveJS           *services.ZWaveJS
}

func newService(conn *ws.WebsocketWriter, ctx context.Context, httpClient *http.HttpClient) *Service {
	return &Service{
		AlarmControlPanel: services.BuildService[services.AlarmControlPanel](conn, ctx),
		Climate:           services.BuildService[services.Climate](conn, ctx),
		Cover:             services.BuildService[services.Cover](conn, ctx),
		Light:             services.BuildService[services.Light](conn, ctx),
		HomeAssistant:     services.BuildService[services.HomeAssistant](conn, ctx),
		Lock:              services.BuildService[services.Lock](conn, ctx),
		MediaPlayer:       services.BuildService[services.MediaPlayer](conn, ctx),
		Switch:            services.BuildService[services.Switch](conn, ctx),
		InputBoolean:      services.BuildService[services.InputBoolean](conn, ctx),
		InputButton:       services.BuildService[services.InputButton](conn, ctx),
		InputText:         services.BuildService[services.InputText](conn, ctx),
		InputDatetime:     services.BuildService[services.InputDatetime](conn, ctx),
		InputNumber:       services.BuildService[services.InputNumber](conn, ctx),
		Event:             services.BuildService[services.Event](conn, ctx),
		Notify:            services.BuildService[services.Notify](conn, ctx),
		Number:            services.BuildService[services.Number](conn, ctx),
		Scene:             services.BuildService[services.Scene](conn, ctx),
		Script:            services.BuildService[services.Script](conn, ctx),
		TTS:               services.BuildService[services.TTS](conn, ctx),
		Vacuum:            services.BuildService[services.Vacuum](conn, ctx),
		ZWaveJS:           services.BuildService[services.ZWaveJS](conn, ctx),
	}
}
