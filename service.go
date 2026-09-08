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

func newService(api services.API) *Service {
	return &Service{
		AdaptiveLighting:  services.BuildService[services.AdaptiveLighting](api),
		AlarmControlPanel: services.BuildService[services.AlarmControlPanel](api),
		Climate:           services.BuildService[services.Climate](api),
		Cover:             services.BuildService[services.Cover](api),
		Light:             services.BuildService[services.Light](api),
		HomeAssistant:     services.BuildService[services.HomeAssistant](api),
		Lock:              services.BuildService[services.Lock](api),
		MediaPlayer:       services.BuildService[services.MediaPlayer](api),
		Switch:            services.BuildService[services.Switch](api),
		InputBoolean:      services.BuildService[services.InputBoolean](api),
		InputButton:       services.BuildService[services.InputButton](api),
		InputText:         services.BuildService[services.InputText](api),
		InputDatetime:     services.BuildService[services.InputDatetime](api),
		InputNumber:       services.BuildService[services.InputNumber](api),
		Event:             services.BuildService[services.Event](api),
		Notify:            services.BuildService[services.Notify](api),
		Number:            services.BuildService[services.Number](api),
		Scene:             services.BuildService[services.Scene](api),
		Script:            services.BuildService[services.Script](api),
		Timer:             services.BuildService[services.Timer](api),
		TTS:               services.BuildService[services.TTS](api),
		Vacuum:            services.BuildService[services.Vacuum](api),
		ZWaveJS:           services.BuildService[services.ZWaveJS](api),
	}
}
