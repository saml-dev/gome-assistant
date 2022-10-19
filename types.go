package gomeassistant

import "time"

type EventZWaveJSValueNotification struct {
	EventType string `yaml:"event_type"`
	Data      struct {
		Domain           string `yaml:"domain"`
		NodeID           int    `yaml:"node_id"`
		HomeID           int64  `yaml:"home_id"`
		Endpoint         int    `yaml:"endpoint"`
		DeviceID         string `yaml:"device_id"`
		CommandClass     int    `yaml:"command_class"`
		CommandClassName string `yaml:"command_class_name"`
		Label            string `yaml:"label"`
		Property         string `yaml:"property"`
		PropertyName     string `yaml:"property_name"`
		PropertyKey      string `yaml:"property_key"`
		PropertyKeyName  string `yaml:"property_key_name"`
		Value            string `yaml:"value"`
		ValueRaw         int    `yaml:"value_raw"`
	} `yaml:"data"`
	Origin    string    `yaml:"origin"`
	TimeFired time.Time `yaml:"time_fired"`
}
