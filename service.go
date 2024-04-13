package gomeassistant

import (
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/internal/websocket"
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

func newService(conn *websocket.Conn, httpClient *http.HttpClient) *Service {
	return &Service{
		AlarmControlPanel: services.BuildService[services.AlarmControlPanel](conn),
		Climate:           services.BuildService[services.Climate](conn),
		Cover:             services.BuildService[services.Cover](conn),
		Light:             services.BuildService[services.Light](conn),
		HomeAssistant:     services.BuildService[services.HomeAssistant](conn),
		Lock:              services.BuildService[services.Lock](conn),
		MediaPlayer:       services.BuildService[services.MediaPlayer](conn),
		Switch:            services.BuildService[services.Switch](conn),
		InputBoolean:      services.BuildService[services.InputBoolean](conn),
		InputButton:       services.BuildService[services.InputButton](conn),
		InputText:         services.BuildService[services.InputText](conn),
		InputDatetime:     services.BuildService[services.InputDatetime](conn),
		InputNumber:       services.BuildService[services.InputNumber](conn),
		Event:             services.BuildService[services.Event](conn),
		Notify:            services.BuildService[services.Notify](conn),
		Number:            services.BuildService[services.Number](conn),
		Scene:             services.BuildService[services.Scene](conn),
		Script:            services.BuildService[services.Script](conn),
		TTS:               services.BuildService[services.TTS](conn),
		Vacuum:            services.BuildService[services.Vacuum](conn),
		ZWaveJS:           services.BuildService[services.ZWaveJS](conn),
	}
}
