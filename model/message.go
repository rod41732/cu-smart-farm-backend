package model

import "encoding/json"

// APICall represent payload for API Calls via WS
type APICall struct {
	EndPoint string                 `json:"endPoint" binding:"required"` // addDevice, removeDevice, setDevice, pollDevice, listDevice...
	Token    string                 `json:"token" binding:"required"`
	Payload  map[string]interface{} `json:"payload" binding:"required"` // json data depend on command
}

type RelayMode struct { // use when set relay mode
	Mode   string      `json:"mode"`
	Detail interface{} `json:"detail"`
}

// AddDeviceMessage is payload format for addDevice API
type AddDeviceMessage struct {
	DeviceID     string `json:"deviceID" binding:"required"`
	DeviceSecret string `json:"deviceSecret" binding:"required"`
}

// RemoveDeviceMessage is payload format for removeDevice API
type RemoveDeviceMessage struct {
	DeviceID string `json:"deviceID" binding:"required"`
}

// FromMap is "constructor" for converting map[string]interface{} to AddDeviceMessage  return error if can't convert
func (message *AddDeviceMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}

// FromMap is "constructor" for converting map[string]interface{} to RemoveDeviceMessage  return error if can't convert
func (message *RemoveDeviceMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}
