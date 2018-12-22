package model

import "encoding/json"

// User : interface type of 'RealUser' and 'NullUser'
type User interface {
	ReportStatus(payload interface{})
}

// Device : basic device info
type Device struct {
	DeviceID     string `json:"id"`
	DeviceSecret string `json:"secret"`
	Owner        string `json:"owner"`
}

// DeviceMessage : mqtt stat from device
type DeviceMessage struct {
	Type string `json:"t" binding:"required"`
	// Payload struct {
	Soil     float32 `json:"Soil" binding:"required"`
	Humidity float32 `json:"Humidity" binding:"required"`
	Temp     float32 `json:"Temp" binding:"required"`
	// Also relay state ON/OFF ??
	// } `json:"data"`
}

// ToMap convert struct back to map
func (dmesg *DeviceMessage) ToMap() (out map[string]interface{}) {
	str, _ := json.Marshal(dmesg)
	json.Unmarshal(str, &out)
	return
}
