package user

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	UserAPI "github.com/rod41732/cu-smart-farm-backend/api/user"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
)

// RealUser represent client connected via WebSocket
type RealUser struct {
	Username     string
	currentToken string
	devices      []string
	conn         *websocket.Conn
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//Init  initializes user
func (user *RealUser) Init(devices []string, conn *websocket.Conn) {
	user.RegenerateToken()
	user.devices = devices
	user.conn = conn
}

// Command set relay state of device (specified via payload)
func (user *RealUser) Command(relay string, workmode string, payload interface{}) {

}

// ReportStatus sends MQTT data to user via WebSocket
func (user *RealUser) ReportStatus(payload interface{}) {
	resp, _ := json.Marshal(payload)
	user.conn.WriteMessage(1, resp) // 1 is text message
}

// AddDevice adds device into user's device list
func (user *RealUser) AddDevice(payload map[string]interface{}) (bool, string) {
	common.Println("addding device ?")
	common.Println(payload)
	str, err := json.Marshal(payload)
	if err != nil {
		return false, "Bad Request"
	}
	var message model.AddDeviceMessage
	json.Unmarshal(str, &message)
	if err != nil {
		return false, "Bad Request"
	}
	return UserAPI.HandleAddDevice(message, user.Username)
}

// RemoveDevice removes device from user's device list
func (user *RealUser) RemoveDevice(payload map[string]interface{}) (bool, string) {
	str, err := json.Marshal(payload)
	if err != nil {
		return false, "Bad Request"
	}
	var message model.RemoveDeviceMessage
	json.Unmarshal(str, &message)
	if err != nil {
		return false, "Bad Request"
	}
	if !(common.StringInSlice(message.DeviceID, user.devices)) {
		return false, "not your device"
	}
	return UserAPI.HandleRemoveDevice(message, user.Username)
}

// RegenerateToken : Regenerate user websocket authorization token
func (user *RealUser) RegenerateToken() string {
	user.currentToken = randomString(20)
	return user.currentToken
}

// CheckToken : Check user websocket authorization token
func (user *RealUser) CheckToken(token string) bool {
	return token == user.currentToken
}

// randomString : helper function for random string with custom length and charset
func randomString(length int) string {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (user *RealUser) ownsDevice(deviceID string) bool {
	return common.StringInSlice(deviceID, user.devices)
}
