package types

type SetTemperatureRequest struct {
	Temperature    float32
	TargetTempHigh float32
	TargetTempLow  float32
	HvacMode       string
}

func (r *SetTemperatureRequest) ToJSON() map[string]any {
	m := map[string]any{}
	if r.Temperature != 0 {
		m["temperature"] = r.Temperature
	}
	if r.TargetTempHigh != 0 {
		m["target_temp_high"] = r.TargetTempHigh
	}
	if r.TargetTempLow != 0 {
		m["target_temp_low"] = r.TargetTempLow
	}
	if r.HvacMode != "" {
		m["hvac_mode"] = r.HvacMode
	}
	return m
}
