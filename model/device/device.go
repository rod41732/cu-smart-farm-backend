package device

import (
	"encoding/json"

	"github.com/rod41732/cu-smart-farm-backend/model"

	"github.com/rod41732/cu-smart-farm-backend/common"
)

// Device : device connected ti
type Device struct {
	ID               string                            `json:"id"`
	Name             string                            `json:"name"`
	Description      string                            `json:"desc"`
	Secret           string                            `json:"secret"`
	Owner            string                            `json:"owner"`
	RelayStates      map[string]RelayState             `json:"state"`
	LastSensorValues model.DeviceMessageV1_0           `json:"-"`
	LastRelays        []int                             `json:"-"`
	PastStates       map[string]map[string]interface{} `json:"pastState"` // store
}

// Version return device's version
// TODO: use actual version
func (device *Device) Version() string {
	return "1.0"
}

// FromMap initialize data using map[string]interface{}
func (device *Device) FromMap(data map[string]interface{}) error {
	str, _ := json.Marshal(data)
	err := json.Unmarshal(str, &device)
	if err != nil {
		return err
	}
	if device.RelayStates == nil {
		device.RelayStates = make(map[string]RelayState)
	}
	if device.PastStates == nil {
		device.PastStates = make(map[string]map[string]interface{})
	}
	for _, key := range common.PossibleRelays {
		if device.PastStates[key] == nil {
			device.PastStates[key] = make(map[string]interface{})
		}
	}
	return nil
}
