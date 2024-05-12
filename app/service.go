package app

import (
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/internal/services"
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

func newService(app *App, httpClient *http.HttpClient) *Service {
	return &Service{
		AlarmControlPanel: services.NewAlarmControlPanel(app),
		Climate:           services.NewClimate(app),
		Cover:             services.NewCover(app),
		Light:             services.NewLight(app),
		HomeAssistant:     services.NewHomeAssistant(app),
		Lock:              services.NewLock(app),
		MediaPlayer:       services.NewMediaPlayer(app),
		Switch:            services.NewSwitch(app),
		InputBoolean:      services.NewInputBoolean(app),
		InputButton:       services.NewInputButton(app),
		InputText:         services.NewInputText(app),
		InputDatetime:     services.NewInputDatetime(app),
		InputNumber:       services.NewInputNumber(app),
		Event:             services.NewEvent(app),
		Notify:            services.NewNotify(app),
		Number:            services.NewNumber(app),
		Scene:             services.NewScene(app),
		Script:            services.NewScript(app),
		TTS:               services.NewTTS(app),
		Vacuum:            services.NewVacuum(app),
		ZWaveJS:           services.NewZWaveJS(app),
	}
}
