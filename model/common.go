package model

import (
	"encoding/json"
)

// DeviceMessagePayload : mqtt message payload (for report device stat) from device
type DeviceMessagePayload struct {
	Soil     float32 `json:"Soil"`
	Humidity float32 `json:"Humidity"`
	Temp     float32 `json:"Temp"`
}

// DeviceMessage : mqtt message from device
type DeviceMessage struct {
	Type    string               `json:"t" binding:"required"`
	Payload DeviceMessagePayload `json:"data"`
}

// // ToMap is convenient method for converting struct back to map
// func (dmesg *DeviceMessage) ToMap() (out map[string]interface{}) {
// 	str, _ := json.Marshal(dmesg)
// 	json.Unmarshal(str, &out)
// 	return
// }

// ToMap is convenient method for converting struct back to map
func (dmesg *DeviceMessagePayload) ToMap() (out map[string]interface{}) {
	str, _ := json.Marshal(dmesg)
	json.Unmarshal(str, &out)
	return
}
