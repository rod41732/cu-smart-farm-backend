package user

import (
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
)

// NullUser is "placeholder" when client is disconnected
type NullUser struct {
}

// ReportStatus : for null user only insert to influxDB
func (user *NullUser) ReportStatus(payload model.DeviceMessage, deviceID string) {
	common.Printf("[NullUser] <<< %v\n", payload)
	out := payload.ToMap()
	delete(out, "t")
	common.WriteInfluxDB("cu_smartfarm_sensor_log", map[string]string{"device": deviceID}, out)
	// return
}
