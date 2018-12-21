package user

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"github.com/rod41732/cu-smart-farm-backend/api/user"
	model "github.com/rod41732/cu-smart-farm-backend/model"
)

// User struct represent an User object which can call API
type User struct {
	username     string
	currentToken string
	conn         *websocket.Conn
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
	return user.HandleAddDevice(message, caller.username)
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
	return user.HandleRemoveDevice(message, caller.username)
}

// New is user constructor
func New(username, currentToken string, conn *websocket.Conn) User {
	return User{username: username, currentToken: currentToken, conn: conn}
}
