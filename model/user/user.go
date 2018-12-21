package user

import (
	"encoding/json"

	model ".."
	"../../api/user"
)

// User struct represent an User object which can call API
type User struct {
	Username     string
	CurrentToken string
}

// type User__ interface {
// 	Command(deviceID string, relay string, workmode string, payload interface{})
// 	ReportStatus(deviceID string, payload interface{})
// 	AddDevice(deviceID string)
// 	RemoveDevice(deviceID string)
// }

// AddDevice adds device into user's device list
func (caller *User) AddDevice(payload map[string]interface{}) (bool, string) {
	str, err := json.Marshal(payload)
	if err != nil {
		return false, "Bad Request"
	}
	var message model.AddDeviceMessage
	json.Unmarshal(str, message)
	if err != nil {
		return false, "Bad Request"
	}
	return user.HandleAddDevice(message, caller.Username)
}

// RemoveDevice removes device from user's device list
func (caller *User) RemoveDevice(payload map[string]interface{}) (bool, string) {
	str, err := json.Marshal(payload)
	if err != nil {
		return false, "Bad Request"
	}
	var message model.RemoveDeviceMessage
	json.Unmarshal(str, message)
	if err != nil {
		return false, "Bad Request"
	}
	return user.HandleRemoveDevice(message, caller.Username)
}
