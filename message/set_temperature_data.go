package message

type SetTemperatureData struct {
	// FIXME: 0° is not an impossible temperature, so treating zero as
	// "missing" here is a little bit dubious.

	Temperature    float32 `json:"temperature,omitempty"`
	TargetTempHigh float32 `json:"target_temp_high,omitempty"`
	TargetTempLow  float32 `json:"target_temp_low,omitempty"`
	HvacMode       string  `json:"hvac_mode,omitempty"`
}
