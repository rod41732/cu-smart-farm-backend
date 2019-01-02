package user

/*
// NullUser is "placeholder" when client is disconnected
type NullUser struct {
}

// ReportStatus : for null user only insert to influxDB
func (user *NullUser) ReportStatus(payload model.DeviceMessagePayload, deviceID string) {
	common.Printf("[NullUser] <<< %v\n", payload)
	common.WriteInfluxDB("cu_smartfarm_sensor_log", map[string]string{"device": deviceID}, payload.ToMap())
	// return
}
*/
