package types

type NotifyRequest struct {
	// Which notify service to call, such as mobile_app_sams_iphone
	ServiceName string
	Message     string
	Title       string
	Data        map[string]any
}

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
