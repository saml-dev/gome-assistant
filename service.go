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

func newService(ctx context.Context, conn *ws.WebsocketWriter, httpClient *http.HttpClient) *Service {
	return &Service{
		AlarmControlPanel: services.BuildService[services.AlarmControlPanel](ctx, conn),
		Climate:           services.BuildService[services.Climate](ctx, conn),
		Cover:             services.BuildService[services.Cover](ctx, conn),
		Light:             services.BuildService[services.Light](ctx, conn),
		HomeAssistant:     services.BuildService[services.HomeAssistant](ctx, conn),
		Lock:              services.BuildService[services.Lock](ctx, conn),
		MediaPlayer:       services.BuildService[services.MediaPlayer](ctx, conn),
		Switch:            services.BuildService[services.Switch](ctx, conn),
		InputBoolean:      services.BuildService[services.InputBoolean](ctx, conn),
		InputButton:       services.BuildService[services.InputButton](ctx, conn),
		InputText:         services.BuildService[services.InputText](ctx, conn),
		InputDatetime:     services.BuildService[services.InputDatetime](ctx, conn),
		InputNumber:       services.BuildService[services.InputNumber](ctx, conn),
		Event:             services.BuildService[services.Event](ctx, conn),
		Notify:            services.BuildService[services.Notify](ctx, conn),
		Number:            services.BuildService[services.Number](ctx, conn),
		Scene:             services.BuildService[services.Scene](ctx, conn),
		Script:            services.BuildService[services.Script](ctx, conn),
		TTS:               services.BuildService[services.TTS](ctx, conn),
		Vacuum:            services.BuildService[services.Vacuum](ctx, conn),
		ZWaveJS:           services.BuildService[services.ZWaveJS](ctx, conn),
	}
}
