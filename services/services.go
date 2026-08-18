// Package services contains the public types used to call Home Assistant
// services.
package services

import internalservices "saml.dev/gome-assistant/internal/services"

// BaseServiceRequest contains the fields needed to call a Home Assistant
// service. ServiceData may contain any JSON-serializable value accepted by the
// service.
type BaseServiceRequest = internalservices.BaseServiceRequest

// Target identifies the entities, areas, or devices affected by a service
// call.
type Target = internalservices.Target

type AdaptiveLighting = internalservices.AdaptiveLighting
type AlarmControlPanel = internalservices.AlarmControlPanel
type Climate = internalservices.Climate
type Cover = internalservices.Cover
type Event = internalservices.Event
type HomeAssistant = internalservices.HomeAssistant
type InputBoolean = internalservices.InputBoolean
type InputButton = internalservices.InputButton
type InputDatetime = internalservices.InputDatetime
type InputNumber = internalservices.InputNumber
type InputText = internalservices.InputText
type Light = internalservices.Light
type Lock = internalservices.Lock
type MediaPlayer = internalservices.MediaPlayer
type Notify = internalservices.Notify
type Number = internalservices.Number
type Scene = internalservices.Scene
type Script = internalservices.Script
type Switch = internalservices.Switch
type Timer = internalservices.Timer
type TTS = internalservices.TTS
type Vacuum = internalservices.Vacuum
type ZWaveJS = internalservices.ZWaveJS

// Entities returns a target containing the provided entity IDs.
func Entities(entityIDs ...string) Target {
	return internalservices.Entities(entityIDs...)
}

// Areas returns a target containing the provided area IDs.
func Areas(areaIDs ...string) Target {
	return internalservices.Areas(areaIDs...)
}

// Devices returns a target containing the provided device IDs.
func Devices(deviceIDs ...string) Target {
	return internalservices.Devices(deviceIDs...)
}
