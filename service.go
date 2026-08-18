package gomeassistant

import (
	internalservices "saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/services"
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
		AdaptiveLighting:  internalservices.BuildService[internalservices.AdaptiveLighting](app),
		AlarmControlPanel: internalservices.BuildService[internalservices.AlarmControlPanel](app),
		Climate:           internalservices.BuildService[internalservices.Climate](app),
		Cover:             internalservices.BuildService[internalservices.Cover](app),
		Light:             internalservices.BuildService[internalservices.Light](app),
		HomeAssistant:     internalservices.BuildService[internalservices.HomeAssistant](app),
		Lock:              internalservices.BuildService[internalservices.Lock](app),
		MediaPlayer:       internalservices.BuildService[internalservices.MediaPlayer](app),
		Switch:            internalservices.BuildService[internalservices.Switch](app),
		InputBoolean:      internalservices.BuildService[internalservices.InputBoolean](app),
		InputButton:       internalservices.BuildService[internalservices.InputButton](app),
		InputText:         internalservices.BuildService[internalservices.InputText](app),
		InputDatetime:     internalservices.BuildService[internalservices.InputDatetime](app),
		InputNumber:       internalservices.BuildService[internalservices.InputNumber](app),
		Event:             internalservices.BuildService[internalservices.Event](app),
		Notify:            internalservices.BuildService[internalservices.Notify](app),
		Number:            internalservices.BuildService[internalservices.Number](app),
		Scene:             internalservices.BuildService[internalservices.Scene](app),
		Script:            internalservices.BuildService[internalservices.Script](app),
		Timer:             internalservices.BuildService[internalservices.Timer](app),
		TTS:               internalservices.BuildService[internalservices.TTS](app),
		Vacuum:            internalservices.BuildService[internalservices.Vacuum](app),
		ZWaveJS:           internalservices.BuildService[internalservices.ZWaveJS](app),
	}
}
