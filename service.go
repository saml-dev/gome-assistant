package gomeassistant

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
	conn := app.wsConn
	return &Service{
		AlarmControlPanel: services.NewAlarmControlPanel(conn),
		Climate:           services.NewClimate(conn),
		Cover:             services.NewCover(conn),
		Light:             services.NewLight(app),
		HomeAssistant:     services.NewHomeAssistant(conn),
		Lock:              services.NewLock(conn),
		MediaPlayer:       services.NewMediaPlayer(conn),
		Switch:            services.NewSwitch(conn),
		InputBoolean:      services.NewInputBoolean(conn),
		InputButton:       services.NewInputButton(conn),
		InputText:         services.NewInputText(conn),
		InputDatetime:     services.NewInputDatetime(conn),
		InputNumber:       services.NewInputNumber(conn),
		Event:             services.NewEvent(app),
		Notify:            services.NewNotify(conn),
		Number:            services.NewNumber(conn),
		Scene:             services.NewScene(conn),
		Script:            services.NewScript(conn),
		TTS:               services.NewTTS(conn),
		Vacuum:            services.NewVacuum(conn),
		ZWaveJS:           services.NewZWaveJS(conn),
	}
}
