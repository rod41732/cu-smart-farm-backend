package device

import "encoding/json"

// Device : device connected ti
type Device struct {
	ID          string                `json:"id" binding:"required"`
	Name        string                `json:"name" binding:"required"`
	Secret      string                `json:"secret" binding:"required"`
	Owner       string                `json:"owner" binding:"required"`
	RelayStates map[string]RelayState `json:"state" binding:"required"`
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
	return nil
}
