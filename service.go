package gomeassistant

import (
	"saml.dev/gome-assistant/internal/services"
)

type Service struct {
	AdaptiveLighting  *services.AdaptiveLighting
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
	Timer             *services.Timer
	TTS               *services.TTS
	Vacuum            *services.Vacuum
	ZWaveJS           *services.ZWaveJS
}

func newService(app *App) *Service {
	return &Service{
		AdaptiveLighting:  services.BuildService[services.AdaptiveLighting](app),
		AlarmControlPanel: services.BuildService[services.AlarmControlPanel](app),
		Climate:           services.BuildService[services.Climate](app),
		Cover:             services.BuildService[services.Cover](app),
		Light:             services.BuildService[services.Light](app),
		HomeAssistant:     services.BuildService[services.HomeAssistant](app),
		Lock:              services.BuildService[services.Lock](app),
		MediaPlayer:       services.BuildService[services.MediaPlayer](app),
		Switch:            services.BuildService[services.Switch](app),
		InputBoolean:      services.BuildService[services.InputBoolean](app),
		InputButton:       services.BuildService[services.InputButton](app),
		InputText:         services.BuildService[services.InputText](app),
		InputDatetime:     services.BuildService[services.InputDatetime](app),
		InputNumber:       services.BuildService[services.InputNumber](app),
		Event:             services.BuildService[services.Event](app),
		Notify:            services.BuildService[services.Notify](app),
		Number:            services.BuildService[services.Number](app),
		Scene:             services.BuildService[services.Scene](app),
		Script:            services.BuildService[services.Script](app),
		Timer:             services.BuildService[services.Timer](app),
		TTS:               services.BuildService[services.TTS](app),
		Vacuum:            services.BuildService[services.Vacuum](app),
		ZWaveJS:           services.BuildService[services.ZWaveJS](app),
	}
}
