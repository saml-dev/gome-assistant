package gomeassistant

import (
	"context"

	"github.com/gorilla/websocket"
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/internal/services"
)

type Service struct {
	HomeAssistant *services.HomeAssistant
	Light         *services.Light
	Lock          *services.Lock
	Switch        *services.Switch
	InputBoolean  *services.InputBoolean
	InputButton   *services.InputButton
	InputText     *services.InputText
	InputDatetime *services.InputDatetime
	InputNumber   *services.InputNumber
	Notify        *services.Notify
	Number        *services.Number
}

func newService(conn *websocket.Conn, ctx context.Context, httpClient *http.HttpClient) *Service {
	return &Service{
		Light:         services.BuildService[services.Light](conn, ctx),
		HomeAssistant: services.BuildService[services.HomeAssistant](conn, ctx),
		Lock:          services.BuildService[services.Lock](conn, ctx),
		Switch:        services.BuildService[services.Switch](conn, ctx),
		InputBoolean:  services.BuildService[services.InputBoolean](conn, ctx),
		InputButton:   services.BuildService[services.InputButton](conn, ctx),
		InputText:     services.BuildService[services.InputText](conn, ctx),
		InputDatetime: services.BuildService[services.InputDatetime](conn, ctx),
		InputNumber:   services.BuildService[services.InputNumber](conn, ctx),
		Notify:        services.BuildService[services.Notify](conn, ctx),
		Number:        services.BuildService[services.Number](conn, ctx),
	}
}
