package app

type ZWaveJSEventData struct {
	Domain           string `json:"domain"`
	NodeID           int    `json:"node_id"`
	HomeID           int64  `json:"home_id"`
	Endpoint         int    `json:"endpoint"`
	DeviceID         string `json:"device_id"`
	CommandClass     int    `json:"command_class"`
	CommandClassName string `json:"command_class_name"`
	Label            string `json:"label"`
	Property         string `json:"property"`
	PropertyName     string `json:"property_name"`
	PropertyKey      string `json:"property_key"`
	PropertyKeyName  string `json:"property_key_name"`
	Value            string `json:"value"`
	ValueRaw         int    `json:"value_raw"`
}
