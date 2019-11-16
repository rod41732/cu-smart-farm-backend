package model

import (
	"fmt"
)
// ensure that DeviceMessageV1_0 implement the interface
var _ DeviceMessageData = &DeviceMessageV1_0{} 

// DeviceMessageV1_0 : mqtt message from device v1.0
type DeviceMessageV1_0 struct {
	Relay [5]int `json:"r"`
	Temp float32 `json:"t"`
	AirHumid float32 `json:"h"`
	SoilHumid float32 `json:"sh"`
}

// Version return 1.0
func (msg *DeviceMessageV1_0) Version() string {
	return "1.0"
}

// ToInflux convert message to influx measurements
func (msg *DeviceMessageV1_0) ToInflux() map[string]interface{} {
	out := make(map[string]interface{})
	for i, val := range msg.Relay {
		out[fmt.Sprintf("Relay%d", i+1)] = val
	}
	out["temp"] = msg.Temp
	out["humid"] = msg.AirHumid
	out["soil"] = msg.SoilHumid
	return out
}

