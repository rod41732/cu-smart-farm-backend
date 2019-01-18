package model

import (
	"encoding/json"
)

// DeviceMessagePayload : mqtt message payload (for report device stat) from device
type DeviceMessagePayload struct {
	Soil     float32 `json:"Soil"`
	Humidity float32 `json:"Humidity"`
	Temp     float32 `json:"Temp"`
	Relay1   string  `json:"Relay1"`
	Relay2   string  `json:"Relay2"`
	Relay3   string  `json:"Relay3"`
	Relay4   string  `json:"Relay4"`
	Relay5   string  `json:"Relay5"`
}

// DeviceMessage : mqtt message from device
type DeviceMessage struct {
	Type    string               `json:"t"`
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
