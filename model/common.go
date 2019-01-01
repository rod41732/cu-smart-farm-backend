package model

import (
	"encoding/json"
)

// User : interface type of 'RealUser' and 'NullUser'
type User interface {
	ReportStatus(payload DeviceMessage, deviceID string)
}

// RelayState represents state of a relay (On/Off) and it's detail

// DeviceMessage : mqtt stat from device
type DeviceMessage struct {
	Type string `json:"t" binding:"required"`
	// Payload struct {
	Soil     float32 `json:"Soil"`
	Humidity float32 `json:"Humidity"`
	Temp     float32 `json:"Temp"`
	// Also relay state ON/OFF ??
	// } `json:"data"`
}

// ToMap is convenient method for converting struct back to map
func (dmesg *DeviceMessage) ToMap() (out map[string]interface{}) {
	str, _ := json.Marshal(dmesg)
	json.Unmarshal(str, &out)
	return
}
