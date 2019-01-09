package device

import (
	"encoding/json"

	"github.com/rod41732/cu-smart-farm-backend/common"
)

// Device : device connected ti
type Device struct {
	ID          string                            `json:"id"`
	Name        string                            `json:"name"`
	Secret      string                            `json:"secret"`
	Owner       string                            `json:"owner"`
	RelayStates map[string]RelayState             `json:"state"`
	PastStates  map[string]map[string]interface{} `json:"pastState"`
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
